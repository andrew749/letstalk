package main

import (
	"flag"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/utility"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

func getTimeXDaysAgo(days int) string {
	now := time.Now()
	res := now.Add(-time.Duration(days*24) * time.Hour)
	return res.Format("2006-01-02")
}

const (
	NUM_DAYS_TILL_NOTIFICATION = 30
)

var (
	runId = flag.String("runId", time.Now().String(), "Id to mark these notifications with")
)

const (
	MENTOR_TITLE   = `Reconnect with your mentee!!`
	MENTOR_MESSAGE = `Its been almost a month since you last connected with
	your mentee. You should reach out to see how they are doing`
	MENTEE_TITLE   = `Reconnect with your mentor!!`
	MENTEE_MESSAGE = `Its been almost a month since you last connected with
		your mentor. You should reach out if you have any questions or
		just want to reconnect!`
)

func main() {
	utility.Bootstrap()

	// find all relationships that are a month old and send out reconnection notifications
	err := utility.RunWithDb(func(tx *gorm.DB) error {
		// find all matchings a month ago
		// TODO: find a more efficient way of doing this.
		// this is going to pull a lot of data in memory
		startTime := getTimeXDaysAgo(NUM_DAYS_TILL_NOTIFICATION)
		rows, err := tx.
			Table("matchings").
			Select("matchings.id, s.id").
			Where("matchings.created_at < ?", startTime).
			Joins("left join sent_monthly_notifications s on s.matching_id=matchings.id").
			Rows()

		if err != nil {
			rlog.Error(err)
			return err
		}

		rlog.Debugf("Finding matchings from before %s", startTime)

		// don't preallocate the whole array size to save some mem
		// only keep entries with null id on sent_monthly_notifications
		toSend := make([]uint, 0)

		type resStruct struct {
			MatchingId         uint
			SendNotificationId *uint
		}
		var tempResStruct resStruct

		for rows.Next() {
			err := rows.Scan(&tempResStruct.MatchingId, &tempResStruct.SendNotificationId)
			if err != nil {
				rlog.Error(err)
				panic(err)
			}

			if tempResStruct.SendNotificationId == nil {
				toSend = append(toSend, tempResStruct.MatchingId)
				rlog.Debugf("Adding Notification for matching with id %d", tempResStruct.MatchingId)
			} else {
				rlog.Debugf("Monthly Notification for matching with id %d already sent", tempResStruct.MatchingId)
			}

		}

		for _, matchingId := range toSend {
			rlog.Debugf("Sending notification for matching with id %d", matchingId)
			// send out notifications for each of these matchings

			// get the matching data object
			matching, err := data.GetMatchingWithId(tx, matchingId)
			if err != nil {
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}

			// get the user data for users involved in the matching
			mentor, err := query.GetUserProfileById(tx, matching.Mentor, false)
			if err != nil {
				rlog.Error(err)
				continue
			}

			mentee, err := query.GetUserProfileById(tx, matching.Mentee, false)
			if err != nil {
				rlog.Error(err)
				continue
			}

			// create notification for mentor
			mentorErr := notifications.CreateAdHocNotification(
				tx,
				matching.Mentor,
				MENTOR_TITLE,
				MENTOR_MESSAGE,
				nil,
				"first_month_notification_mentor.html",
				map[string]interface{}{
					"name": mentee.FirstName,
				},
				nil,
			)

			// create notification for mentee
			menteeErr := notifications.CreateAdHocNotification(
				tx,
				matching.Mentor,
				MENTEE_TITLE,
				MENTEE_MESSAGE,
				nil,
				"first_month_notification_mentee.html",
				map[string]interface{}{
					"name": mentor.FirstName,
				},
				nil,
			)
			// in reality, either both of these should fail or neither since we're
			// just putting in an sqs queue.
			// For now, not modelling with more granular control (i.e. was the mentor/mentee notification sent)
			if mentorErr != nil {
				rlog.Error(mentorErr)
			}
			if menteeErr != nil {
				rlog.Error(menteeErr)
			}

			// if not error, mark this notification as being sent by a batch job.
			if err := data.SentOutMonthlyNotification(tx, matchingId, *runId); err != nil {
				rlog.Error(err)
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
