package main

import (
	"flag"
	"fmt"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/utility"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

func getTimeXDaysAgo(days int) string {
	now := time.Now()
	res := now.Add(time.Duration(days*24) * time.Hour)
	return res.Format("%x")
}

const (
	NUM_DAYS_TILL_NOTIFICATION = 30
)

var (
	runId = flag.String("runId", time.Now().String(), "Id to mark these notifications with")
)

func main() {
	flag.Parse()
	// find all relationships that are a month old and send out reconnection notifications
	err := utility.RunWithDb(func(tx *gorm.DB) error {
		// find all matchings a month ago
		// TODO: find a more efficient way of doing this.
		// this is going to pull a lot of data in memory
		rows, err := tx.
			Select("matchings.id, s.id").
			Where("created_at < ?", getTimeXDaysAgo(NUM_DAYS_TILL_NOTIFICATION)).
			Joins("left join sent_monthly_notifications s on s.matching_id=id").
			Rows()

		if err != nil {
			rlog.Error(err)
			return err
		}

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

			if tempResStruct.SendNotificationId != nil {
				toSend = append(toSend, tempResStruct.MatchingId)
				rlog.Debug("Adding Notification for matching with id %d", tempResStruct.MatchingId)
			} else {
				rlog.Debug("Monthly Notification for matching with id %d already sent", tempResStruct.MatchingId)
			}

		}

		for _, matchingId := range toSend {
			// send out notifications for each of these
			matching, err := data.GetMatchingWithId(tx, matchingId)
			if err != nil {
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}
			mentorMessage := fmt.Sprintf("Its been almost a month since you last connected with your mentee. You should reach out to see how they are doing")
			menteeMessage := fmt.Sprint("Its been almost a month since you last connected with your mentor. You should reach out if you have any questions or just want to reconnect!")

			mentorErr := notifications.CreateAdHocNotification(
				tx,
				matching.Mentor,
				"Reconnect with your mentee!!",
				mentorMessage,
				nil,
				"first_month_notification_mentor.html",
				map[string]string{
					"name": "Andrew",
				},
			)

			menteeErr := notifications.CreateAdHocNotification(
				tx,
				matching.Mentor,
				"Reconnect with your mentor!!",
				menteeMessage,
				nil,
				"first_month_notification_mentee.html",
				map[string]string{
					"name": "Wojtek",
				},
			)
			if mentorErr != nil {
				rlog.Error(mentorErr)
			}
			if menteeErr != nil {
				rlog.Error(menteeErr)
			}
			// if not error
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
