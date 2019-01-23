package remind_meetup_job

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"letstalk/server/core/ctx"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const REMIND_MEETUP_JOB jobmine.JobType = "RemindMeetupJob"

type MatchType string

const (
	MATCH_TYPE_MENTOR     MatchType = "MENTOR"
	MATCH_TYPE_MENTEE     MatchType = "MENTEE"
	MATCH_TYPE_CONNECTION MatchType = "CONNECTION"
)

const (
	REMINDER_ID_METADATA_KEY   = "reminderId"
	REMINDER_TYPE_METADATA_KEY = "reminderType"
	USER_ID_METADATA_KEY       = "userId"
	MATCH_USER_ID_METADATA_KEY = "matchUserId"
)

func packageTaskRecordMetadata(
	reminderId uint,
	userId data.TUserID,
	matchUserId data.TUserID,
) map[string]interface{} {
	return map[string]interface{}{
		REMINDER_ID_METADATA_KEY:   reminderId,
		USER_ID_METADATA_KEY:       userId,
		MATCH_USER_ID_METADATA_KEY: matchUserId,
	}
}

func parseTaskRecord(taskRecord jobmine.TaskRecord) (reminderId uint, reminderType data.MeetupType, userId data.TUserID, matchUserId data.TUserID) {
	reminderId = taskRecord.Metadata[REMINDER_ID_METADATA_KEY].(uint)
	reminderType = data.MeetupType(taskRecord.Metadata[REMINDER_TYPE_METADATA_KEY].(string))
	userId = data.TUserID(taskRecord.Metadata[USER_ID_METADATA_KEY].(uint))
	matchUserId = data.TUserID(taskRecord.Metadata[MATCH_USER_ID_METADATA_KEY].(uint))
	return reminderId, reminderType, userId, matchUserId
}

func getMatchType(userId data.TUserID, connection *data.Connection) MatchType {
	if connection.Mentorship == nil {
		return MATCH_TYPE_CONNECTION
	} else if connection.Mentorship.MentorUserId == userId {
		return MATCH_TYPE_MENTEE
	} else {
		return MATCH_TYPE_MENTOR
	}
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	reminderId, reminderType, userId, matchUserId := parseTaskRecord(taskRecord)
	connection, err := query.GetConnectionDetailsUndirected(db, userId, matchUserId)
	if err != nil {
		return nil, err
	}
	if connection == nil {
		return nil, errors.New(fmt.Sprintf("Meetup reminder failed to find connection for users (%d, %d)", userId, matchUserId))
	}

	var user, matchUser *data.User
	if connection.UserOneId == userId {
		user = connection.UserOne
		matchUser = connection.UserTwo
	} else {
		user = connection.UserTwo
		matchUser = connection.UserOne
	}

	matchType := getMatchType(userId, connection)
	templateParams := taskRecord.Metadata

	// TODO(aklen): create updated meetup reminder template
	dbErr := ctx.WithinTx(db, func(tx *gorm.DB) error {
		// Mark reminder as processed.
		if err := tx.Model(&data.MeetupReminder{}).
			Update(&data.MeetupReminder{Model: gorm.Model{ID: reminderId}, State: data.MEETUP_REMINDER_PROCESSED}).
			Error; err != nil {
			return err
		}
		// Automatically schedule backup notification in three days.
		backup := &data.MeetupReminder{
			UserId:      userId,
			MatchUserId: matchUserId,
			Type:        reminderType,
			State:       data.MEETUP_REMINDER_SCHEDULED,
			ScheduledAt: time.Now().AddDate(0, 0, 3),
		}
		if err := tx.Model(&data.MeetupReminder{}).Create(&backup).Error; err != nil {
			return err
		}
		if err := notifications.CreateAdHocNotificationNoTransaction(
			tx,
			user.UserId,
			"Reminder to Meet Up",
			fmt.Sprintf("Meet up with your %s, %s!", strings.ToLower(string(matchType)), matchUser.FirstName),
			nil,
			"remind_meetup_notification.html",
			templateParams,
			&jobRecord.RunId); err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return nil, dbErr
	}
	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	reminderId, reminderType, userId, _ := parseTaskRecord(taskRecord)
	rlog.Infof("Unable to send reminder %d (%s) to user with id=%d: %+v", reminderId, reminderType, userId, err)
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	reminderId, reminderType, userId, _ := parseTaskRecord(taskRecord)
	rlog.Infof("Successfully sent reminder %d (%s) to user with id=%d", reminderId, reminderType, userId)
}

var reminderTaskSpec = jobmine.TaskSpec{
	Execute:   execute,
	OnError:   onError,
	OnSuccess: onSuccess,
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	reminders, err := query.GetMeetupRemindersScheduledBefore(db, time.Now())
	if err != nil {
		return nil, err
	}
	metadata := make([]jobmine.Metadata, 0, len(reminders))
	for _, reminder := range reminders {
		metadata = append(metadata,
			jobmine.Metadata(packageTaskRecordMetadata(
				reminder.ID,
				reminder.UserId,
				reminder.MatchUserId,
			)))
	}
	return metadata, nil
}

var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          REMIND_MEETUP_JOB,
	TaskSpec:         reminderTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

// CreateReminderJob Creates a reminder job record to get run at some point.
func CreateReminderJob(db *gorm.DB, runId string, startTime *time.Time, endTime *time.Time) error {
	metadata := map[string]interface{}{}

	if err := db.Create(&jobmine.JobRecord{
		JobType:  REMIND_MEETUP_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
