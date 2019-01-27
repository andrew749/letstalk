package welcome_back_email_job

import (
	"fmt"
	"strconv"
	"time"

	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/core/verify_link"
	"letstalk/server/data"
	"letstalk/server/email"
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const WELCOME_BACK_EMAIL_JOB jobmine.JobType = "WelcomeBackEmailJob"

type UserType string

const (
	USER_TYPE_MENTOR UserType = "MENTOR"
	USER_TYPE_MENTEE UserType = "MENTEE"
)

//  Task record
const USER_ID_METADATA_KEY = "userId"

//  Job record
const (
	START_TIME_METADATA_KEY   = "startTime"
	END_TIME_METADATA_KEY     = "endTime"
	TEST_USER_ID_METADATA_KEY = "testUserId"
)

func packageTaskRecordMetadata(userId data.TUserID) map[string]interface{} {
	return map[string]interface{}{USER_ID_METADATA_KEY: userId}
}

func parseUserId(userIdIntf interface{}) (*data.TUserID, error) {
	userIdFloat, ok := userIdIntf.(float64)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid userId %v", userIdIntf))
	}
	userId := data.TUserID(uint(userIdFloat))
	return &userId, nil
}

func parseUserInfo(taskRecord jobmine.TaskRecord) (*data.TUserID, error) {
	userIdIntf, ok := taskRecord.Metadata[USER_ID_METADATA_KEY]
	if !ok {
		return nil, errors.New("Task missing userId")
	}
	return parseUserId(userIdIntf)
}

func execute(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
) (interface{}, error) {
	userIdPtr, err := parseUserInfo(taskRecord)
	if err != nil {
		return nil, err
	}
	userId := *userIdPtr
	user, err := query.GetUserById(db, userId)
	if err != nil {
		return nil, err
	}

	linkIdPtr, err := verify_link.CreateLink(
		db,
		userId,
		verify_link.LINK_TYPE_WHITELIST_WINTER_2019,
		nil,
	)
	if err != nil {
		return nil, err
	}
	linkId := *linkIdPtr

	verifyLinkHrefLink := fmt.Sprintf(
		"%s/welcome_back.html?requestId=%s",
		utility.BaseUrl,
		linkId,
	)

	to := mail.NewEmail(user.FirstName, user.Email)
	if err := email.SendWelcomeBackEmail(to, verifyLinkHrefLink); err != nil {
		return nil, err
	}

	return linkId, nil
}

func onError(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
	userIdPtr, parseErr := parseUserInfo(taskRecord)
	if parseErr != nil {
		rlog.Infof(
			"Unable to send email: %+v - couldn't parse userId from task record (%+v)",
			err,
			parseErr,
		)
	} else {
		userId := *userIdPtr
		rlog.Infof("Unable to send email to user with id=%d: %+v", userId, err)
	}
}

func onSuccess(
	db *gorm.DB,
	jobRecord jobmine.JobRecord,
	taskRecord jobmine.TaskRecord,
	res interface{},
) {
	userIdPtr, err := parseUserInfo(taskRecord)
	if err != nil {
		rlog.Infof("Successfully sent email - couldn't parse userId from task record (%+v)", err)
		return
	}
	userId := *userIdPtr
	linkId, ok := res.(data.TVerifyLinkID)
	if !ok {
		rlog.Infof(
			"Successfully sent email to user with id=%d - couldn't get linkId from task record",
			userId,
		)
		return
	}
	rlog.Infof("Successfully sent email to user with id=%d and linkId=%s", userId, linkId)
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

func getTasksToCreate(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
	if testUserIdIntf, ok := jobRecord.Metadata[TEST_USER_ID_METADATA_KEY]; ok {
		testUserIdFloat, ok := testUserIdIntf.(float64)
		if !ok {
			return nil, errors.New("Couldn't parse testUserId")
		}
		testUserId := data.TUserID(uint(testUserIdFloat))
		rlog.Warn(fmt.Sprintf("Only using test user %d", testUserId))
		return []jobmine.Metadata{packageTaskRecordMetadata(testUserId)}, nil
	}

	startTime, err := jobmine.TimeFromJobRecord(jobRecord, START_TIME_METADATA_KEY)
	if err != nil {
		return nil, err
	}
	endTime, err := jobmine.TimeFromJobRecord(jobRecord, END_TIME_METADATA_KEY)
	if err != nil {
		return nil, err
	}

	if startTime == nil {
		rlog.Warn("No start time provided. Finding users from beginning")
	}
	if endTime == nil {
		rlog.Warn("No end time provided. Finding users from beginning")
	}

	users, err := query.GetUsersByCreatedAt(db, startTime, endTime)
	if err != nil {
		return nil, err
	}

	metadatas := make([]jobmine.Metadata, 0)
	for _, user := range users {
		metadata := jobmine.Metadata(packageTaskRecordMetadata(user.UserId))
		metadatas = append(metadatas, metadata)
	}
	return metadatas, nil
}

var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType:          WELCOME_BACK_EMAIL_JOB,
	TaskSpec:         reminderTaskSpec,
	GetTasksToCreate: getTasksToCreate,
}

// CreateReminderJob Creates a reminder job record to get run at some point.
func CreateEmailJob(
	db *gorm.DB,
	runId string,
	startTime *time.Time,
	endTime *time.Time,
	testUserId *data.TUserID,
) error {
	metadata := map[string]interface{}{}
	if startTime != nil {
		metadata[START_TIME_METADATA_KEY] = jobmine.FormatTime(*startTime)
	}
	if endTime != nil {
		metadata[END_TIME_METADATA_KEY] = jobmine.FormatTime(*endTime)
	}
	if testUserId != nil {
		metadata[TEST_USER_ID_METADATA_KEY] = strconv.Itoa(int(*testUserId))
	}

	if err := db.Create(&jobmine.JobRecord{
		JobType:  WELCOME_BACK_EMAIL_JOB,
		RunId:    runId,
		Metadata: metadata,
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
