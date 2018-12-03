package remind_meetup_job

import (
	"errors"
	"fmt"
	"time"

	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const REMIND_MEETUP_JOB jobmine.JobType = "RemindMeetupJob"

type UserType string

const (
	USER_TYPE_MENTOR UserType = "MENTOR"
	USER_TYPE_MENTEE UserType = "MENTEE"
)

const (
	USER_TYPE_METADATA_KEY              = "userType"
	USER_ID_METADATA_KEY                = "userId"
	CONNECTION_USER_ID_METADATA_KEY     = "connectionUserId"
	CONNECTION_FIRST_NAME_METADATA_KEY  = "connectionFirstName"
	CONNECTION_LAST_NAME_METADATA_KEY   = "connectionLastName"
	CONNECTION_PROFILE_PIC_METADATA_KEY = "connectionProfilePic"
)

const (
	START_TIME_METADATA_KEY = "startTime"
	END_TIME_METADATA_KEY   = "endTime"
)

func packageTaskRecordMetadata(
	userId data.TUserID,
	userType UserType,
	connectionUserId data.TUserID,
	connectionFirstName string,
	connectionLastName string,
	connectionProfilePic *string,
) map[string]interface{} {
	return map[string]interface{}{
		USER_ID_METADATA_KEY:                userId,
		USER_TYPE_METADATA_KEY:              userType,
		CONNECTION_USER_ID_METADATA_KEY:     connectionUserId,
		CONNECTION_FIRST_NAME_METADATA_KEY:  connectionFirstName,
		CONNECTION_LAST_NAME_METADATA_KEY:   connectionLastName,
		CONNECTION_PROFILE_PIC_METADATA_KEY: connectionProfilePic,
	}
}

func parseUserInfo(taskRecord jobmine.TaskRecord) (data.TUserID, UserType) {
	userId := data.TUserID(taskRecord.Metadata[USER_ID_METADATA_KEY].(uint))
	userType := UserType(taskRecord.Metadata[USER_TYPE_METADATA_KEY].(string))
	return userId, userType
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	userId, _ := parseUserInfo(taskRecord)
	templateParams := taskRecord.Metadata
	connectionFirstName := templateParams[CONNECTION_FIRST_NAME_METADATA_KEY].(string)

	err := notifications.CreateAdHocNotificationNoTransaction(
		db,
		userId,
		"Reminder to Meet Up",
		fmt.Sprintf("Meet up with %s!", connectionFirstName),
		nil,
		"remind_meetup_notification.html",
		templateParams,
		&jobRecord.RunId,
	)
	if err != nil {
		return nil, err
	}
	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	userId, userType := parseUserInfo(taskRecord)
	rlog.Infof("Unable to send message of type=%s to user with id=%d: %+v", userType, userId, err)
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	userId, userType := parseUserInfo(taskRecord)
	rlog.Infof("Successfully sent message of type=%s to user with id=%d", userType, userId)
}

var reminderTaskSpec = jobmine.TaskSpec{
	Execute:   execute,
	OnError:   onError,
	OnSuccess: onSuccess,
}

func getUserType(userId data.TUserID, mentorUserId data.TUserID) UserType {
	if userId == mentorUserId {
		return USER_TYPE_MENTOR
	} else {
		return USER_TYPE_MENTEE
	}
}

func getTime(jobRecord jobmine.JobRecord, key string) *time.Time {
	if val, ok := jobRecord.Metadata[key]; ok {
		if tme, ok := val.(time.Time); ok {
			return &tme
		}
	}
	return nil
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	startTime := getTime(jobRecord, START_TIME_METADATA_KEY)
	endTime := getTime(jobRecord, END_TIME_METADATA_KEY)

	if startTime == nil {
		rlog.Warn("No start time provided. Finding mentorships from beginning")
	}
	if endTime == nil {
		rlog.Warn("No end time provided. Finding mentorships from beginning")
	}

	connections, err := query.GetMentorshipConnectionsByDate(db, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// TODO(wojtek): Filter out connections that have already received this notification

	metadata := make([]jobmine.Metadata, 0)
	for _, connection := range connections {
		if connection.Mentorship == nil || connection.UserOne == nil || connection.UserTwo == nil {
			return nil, errors.New(fmt.Sprintf("Connection %d is missing data", connection.ConnectionId))
		}

		metadata1 := jobmine.Metadata(packageTaskRecordMetadata(
			connection.UserOneId,
			getUserType(connection.UserOneId, connection.Mentorship.MentorUserId),
			connection.UserTwoId,
			connection.UserTwo.FirstName,
			connection.UserTwo.LastName,
			connection.UserTwo.ProfilePic,
		))
		metadata2 := jobmine.Metadata(packageTaskRecordMetadata(
			connection.UserTwoId,
			getUserType(connection.UserTwoId, connection.Mentorship.MentorUserId),
			connection.UserOneId,
			connection.UserOne.FirstName,
			connection.UserOne.LastName,
			connection.UserOne.ProfilePic,
		))
		metadata = append(metadata, metadata1, metadata2)
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
	if startTime != nil {
		metadata[START_TIME_METADATA_KEY] = *startTime
	}
	if endTime != nil {
		metadata[END_TIME_METADATA_KEY] = *endTime
	}

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
