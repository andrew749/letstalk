package seed_mentorships_job

import (
	"errors"
	"fmt"
	"time"

	"letstalk/server/data"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const SEED_MENTORSHIPS_JOB jobmine.JobType = "SeedMentorshipsJob"

const (
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
	connectionUserId data.TUserID,
	connectionFirstName string,
	connectionLastName string,
	connectionProfilePic *string,
) map[string]interface{} {
	return map[string]interface{}{
		USER_ID_METADATA_KEY:                userId,
		CONNECTION_USER_ID_METADATA_KEY:     connectionUserId,
		CONNECTION_FIRST_NAME_METADATA_KEY:  connectionFirstName,
		CONNECTION_LAST_NAME_METADATA_KEY:   connectionLastName,
		CONNECTION_PROFILE_PIC_METADATA_KEY: connectionProfilePic,
	}
}

func parseUserInfo(taskRecord jobmine.TaskRecord) data.TUserID {
	userId := data.TUserID(uint(taskRecord.Metadata[USER_ID_METADATA_KEY].(float64)))
	return userId
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	return "Success", nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	userId := parseUserInfo(taskRecord)
	rlog.Infof("Unable to send message to user with id=%d: %+v", userId, err)
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	userId := parseUserInfo(taskRecord)
	rlog.Infof("Successfully sent message to user with id=%d", userId)
}

var reminderTaskSpec = jobmine.TaskSpec{
	Execute:   execute,
	OnError:   onError,
	OnSuccess: onSuccess,
}

const TIME_LAYOUT = "2006-01-02T15:04:05Z"

func getTime(jobRecord jobmine.JobRecord, key string) (*time.Time, error) {
	if val, ok := jobRecord.Metadata[key]; ok {
		var (
			timeStr string
			ok      bool
		)
		if timeStr, ok = val.(string); !ok {
			return nil, errors.New(fmt.Sprintf("%s must be a time string", key))
		}
		time, err := time.Parse(TIME_LAYOUT, timeStr)
		if err != nil {
			return nil, err
		}
		return &time, nil
	}
	return nil, nil
}

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	return []jobmine.Metadata{}, nil
}

var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          SEED_MENTORSHIPS_JOB,
	TaskSpec:         reminderTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

// CreateReminderJob Creates a reminder job record to get run at some point.
func CreateReminderJob(db *gorm.DB, runId string, startTime *time.Time, endTime *time.Time) error {
	metadata := map[string]interface{}{}
	if startTime != nil {
		metadata[START_TIME_METADATA_KEY] = startTime.Format(TIME_LAYOUT)
	}
	if endTime != nil {
		metadata[END_TIME_METADATA_KEY] = endTime.Format(TIME_LAYOUT)
	}

	if err := db.Create(&jobmine.JobRecord{
		JobType:  SEED_MENTORSHIPS_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
