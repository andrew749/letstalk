package remind_meetup_job

import (
	"fmt"
	"strings"
	"time"

	"letstalk/server/core/errs"
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

// Task metadata keys
const (
	REMINDER_ID_METADATA_KEY   = "reminderId"
	MEETUP_TYPE_METADATA_KEY   = "meetupType"
	USER_ID_METADATA_KEY       = "userId"
	MATCH_USER_ID_METADATA_KEY = "matchUserId"
)

// Notification template keys
const (
	MEETUP_TYPE_KEY       = "meetupType"
	MATCH_TYPE_KEY        = "matchType"
	MATCH_USER_ID_KEY     = "matchUserId"
	MATCH_FIRST_NAME_KEY  = "matchFirstName"
	MATCH_LAST_NAME_KEY   = "matchLastName"
	MATCH_PROFILE_PIC_KEY = "matchProfilePic"
)

func packageTaskRecordMetadata(
	reminderId uint,
	meetupType data.MeetupType,
	userId data.TUserID,
	matchUserId data.TUserID,
) map[string]interface{} {
	return map[string]interface{}{
		REMINDER_ID_METADATA_KEY:   reminderId,
		MEETUP_TYPE_METADATA_KEY:   meetupType,
		USER_ID_METADATA_KEY:       userId,
		MATCH_USER_ID_METADATA_KEY: matchUserId,
	}
}

func packageNotificationData(matchType MatchType, meetupType data.MeetupType, matchUser *data.User) map[string]interface{} {
	return map[string]interface{}{
		MATCH_TYPE_KEY:        matchType,
		MEETUP_TYPE_KEY:       meetupType,
		MATCH_USER_ID_KEY:     matchUser.UserId,
		MATCH_FIRST_NAME_KEY:  matchUser.FirstName,
		MATCH_LAST_NAME_KEY:   matchUser.LastName,
		MATCH_PROFILE_PIC_KEY: matchUser.ProfilePic,
	}
}

func parseTaskRecord(taskRecord jobmine.TaskRecord) (reminderId uint, meetupType data.MeetupType, userId data.TUserID, matchUserId data.TUserID) {
	reminderId = uint(taskRecord.Metadata[REMINDER_ID_METADATA_KEY].(float64))
	meetupType = data.MeetupType(taskRecord.Metadata[MEETUP_TYPE_METADATA_KEY].(string))
	userId = data.TUserID(uint(taskRecord.Metadata[USER_ID_METADATA_KEY].(float64)))
	matchUserId = data.TUserID(uint(taskRecord.Metadata[MATCH_USER_ID_METADATA_KEY].(float64)))
	return reminderId, meetupType, userId, matchUserId
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

func markReminderProcessed(tx *gorm.DB, meetupReminderId uint) error {
	if err := tx.Model(&data.MeetupReminder{}).
		Update(&data.MeetupReminder{Model: gorm.Model{ID: meetupReminderId}, State: data.MEETUP_REMINDER_SENT}).
		Error; err != nil {
		return err
	}
	return nil
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	reminderId, meetupType, userId, matchUserId := parseTaskRecord(taskRecord)
	connection, err := query.GetConnectionDetailsUndirected(db, userId, matchUserId)
	if err != nil {
		return nil, err
	}
	if connection == nil {
		rlog.Errorf("Meetup reminder failed to find connection for users (%d, %d), not reprocessing", userId, matchUserId)
		if err := markReminderProcessed(db, reminderId); err != nil {
			return nil, err
		}
		return nil, errs.NewBaseError("Meetup reminder failed to find connection for users (%d, %d)", userId, matchUserId)
	}

	matchUser, err := query.GetUserById(db, matchUserId)
	if err != nil {
		return nil, err
	}

	matchType := getMatchType(userId, connection)
	templateParams := packageNotificationData(matchType, meetupType, matchUser)

	if err := markReminderProcessed(db, reminderId); err != nil {
		return nil, err
	}
	// Automatically schedule backup notification in three days.
	backup := &data.MeetupReminder{
		UserId:      userId,
		MatchUserId: matchUserId,
		Type:        meetupType,
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: time.Now().AddDate(0, 0, 3),
	}
	if err := db.Model(&data.MeetupReminder{}).Create(&backup).Error; err != nil {
		return nil, err
	}
	rlog.Info("Creating meetup notification with params: %v", templateParams)
	if err := notifications.CreateAdHocNotificationNoTransaction(
		db,
		userId,
		"Reminder to Meet Up",
		fmt.Sprintf("Meet up with your %s, %s!", strings.ToLower(string(matchType)), matchUser.FirstName),
		nil,
		"remind_meetup_notification.html",
		templateParams,
		&jobRecord.RunId); err != nil {
		return nil, err
	}
	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	reminderId, meetupType, userId, _ := parseTaskRecord(taskRecord)
	rlog.Infof("Unable to send reminder %d (%s) to user with id=%d: %+v", reminderId, meetupType, userId, err)
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	reminderId, meetupType, userId, _ := parseTaskRecord(taskRecord)
	rlog.Infof("Successfully sent reminder %d (%s) to user with id=%d", reminderId, meetupType, userId)
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
				reminder.Type,
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
func CreateReminderJob(db *gorm.DB, runId string) error {
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
