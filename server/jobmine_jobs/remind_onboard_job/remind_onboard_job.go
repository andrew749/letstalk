package remind_onboard_job

import (
	"fmt"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const RemindOnboardJob jobmine.JobType = "RemindOnboardJob"

type ReminderType string

const (
	REMINDER_TYPE_TRAIT    ReminderType = "REMINDER_TYPE_TRAIT"
	REMINDER_TYPE_BIO                   = "REMINDER_TYPE_BIO"
	REMINDER_TYPE_POSITION              = "REMINDER_TYPE_POSITION"
	REMINDER_TYPE_GROUP                 = "REMINDER_TYPE_GROUP"
)

const (
	ReminderTypeMetadataKey = "reminderType"
	UserIdMetadataKey       = "userId"
)

// extract properties from map
func parseTaskRecordMetadata(taskRecord jobmine.TaskRecord) (data.TUserID, ReminderType) {
	return data.TUserID(taskRecord.Metadata[UserIdMetadataKey].(float64)), ReminderType(taskRecord.Metadata[ReminderTypeMetadataKey].(string))
}

// take properties and put into map
func packageTaskRecordMetadata(userId data.TUserID, reminderType ReminderType) map[string]interface{} {
	return map[string]interface{}{
		UserIdMetadataKey:       userId,
		ReminderTypeMetadataKey: reminderType,
	}
}

var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType: RemindOnboardJob,
	TaskSpec: jobmine.TaskSpec{
		Execute: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord) (interface{}, error) {
			rlog.Infof("Running Task with id %s", taskRecord.ID)

			// send messages
			userId, notificationType := parseTaskRecordMetadata(taskRecord)
			var (
				templatePath string
				title        string
				message      string
			)

			switch notificationType {
			case REMINDER_TYPE_TRAIT:
				templatePath = ""
				title = ""
				message = ""
				break
			case REMINDER_TYPE_BIO:
				templatePath = ""
				title = ""
				message = ""
				break
			case REMINDER_TYPE_POSITION:
				templatePath = ""
				title = ""
				message = ""
				break
			case REMINDER_TYPE_GROUP:
				templatePath = ""
				title = ""
				message = ""
				break
			}
			if err := notifications.CreateAdHocNotificationNoTransaction(db, userId, title, message, nil, templatePath, map[string]interface{}{}, &jobRecord.RunId); err != nil {
				rlog.Errorf("Unable to create notification for user %d because %+v", userId, err)
				return nil, err
			}
			return "Success", nil
		},
		OnError: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
			userId, notificationType := parseTaskRecordMetadata(taskRecord)
			rlog.Infof("Unable to send message of type=%s to user with id=%d: %+v", notificationType, userId, err)
		},
		OnSuccess: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, res interface{}) {
			userId, notificationType := parseTaskRecordMetadata(taskRecord)
			rlog.Infof("Successfully sent message of type=%s to user with id=%d", notificationType, userId)
		},
	},
	GetTasksToCreate: func(db *gorm.DB, jobRecord jobmine.JobRecord) (*[]jobmine.Metadata, error) {
		// keep track of which users we have sent to so far
		seenUsersSet := make(map[data.TUserID]bool)
		seenUsers := make([]data.TUserID, 0)

		// find all users who haven't entered traits
		usersWithoutTraits, err := usersWithoutTraits(db)
		if err != nil {
			rlog.Errorf("Unable to get users without traits: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without traits %v", len(*usersWithoutTraits), *usersWithoutTraits)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutTraits, seenUsers, seenUsersSet)
		metadata := createReminderNotificationPayloadMetadataForUsers(
			*usersWithoutTraits,
			REMINDER_TYPE_TRAIT,
		)

		// find all users who haven't filled in bio
		usersWithoutBio, err := usersWithoutBio(db, seenUsers)
		if err != nil {
			rlog.Errorf("Unable to get users without bio: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without bio %v", len(*usersWithoutBio), *usersWithoutBio)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutBio, seenUsers, seenUsersSet)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(*usersWithoutBio, REMINDER_TYPE_BIO)...,
		)

		// find all users who havent put in a position
		usersWithoutPosition, err := usersWithoutPosition(db, seenUsers)
		if err != nil {
			rlog.Errorf("Unable to get users without position: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without position %v", len(*usersWithoutPosition), *usersWithoutPosition)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutPosition, seenUsers, seenUsersSet)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(*usersWithoutPosition, REMINDER_TYPE_POSITION)...,
		)

		// find all users who haven't put in group
		usersWithoutGroup, err := usersWithoutGroup(db, seenUsers)
		if err != nil {
			rlog.Errorf("Unable to get users without group: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without group %v", len(*usersWithoutGroup), *usersWithoutGroup)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutGroup, seenUsers, seenUsersSet)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(*usersWithoutGroup, REMINDER_TYPE_GROUP)...,
		)

		// convert map to slice
		res := make([]jobmine.Metadata, 0, len(metadata))
		for _, metadata := range metadata {
			res = append(res, metadata)
		}
		rlog.Debugf("Got %d total updates: %+v", len(res), res)
		return &res, nil
	},
}

func mergeMarkSeen(users []data.TUserID, seenList []data.TUserID, seenSet map[data.TUserID]bool) ([]data.TUserID, map[data.TUserID]bool) {
	for _, user := range users {
		// if this user has yet to be seen add them to the set and list
		if _, ok := seenSet[user]; !ok {
			seenSet[user] = true
			seenList = append(seenList, user)
		}
	}
	return seenList, seenSet
}

func createReminderNotificationPayloadMetadataForUsers(
	users []data.TUserID,
	reminderType ReminderType,
) []jobmine.Metadata {
	res := make([]jobmine.Metadata, 0, len(users))
	for _, user := range users {
		res = append(res, jobmine.Metadata(packageTaskRecordMetadata(user, reminderType)))
	}
	return res
}

func usersWithoutTraits(db *gorm.DB) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Pluck("user_id", &temp).
		Joins("left join user_simple_traits as traits on traits.user_id = users.user_id").
		Having("traits.id = NULL").
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutBio(db *gorm.DB, seenSoFar []data.TUserID) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(data.UserAdditionalData{}).
		Pluck("user_id", &temp).
		Where("bio=NULL").
		Where("user_id not in (?)", seenSoFar).
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutPosition(db *gorm.DB, seenSoFar []data.TUserID) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Pluck("user_id", &temp).
		Where("user_id not in (?)", seenSoFar).
		Joins("left join user_positions as positions on positions.user_id = users.user_id").
		Having("positions.id is NULL").
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutGroup(db *gorm.DB, seenSoFar []data.TUserID) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Pluck("user_id", &temp).
		Where("user_id not in (?)", seenSoFar).
		Joins("left join user_groups as groups on groups.user_id = users.user_id").
		Having("groups.id is NULL").
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

// CreateReminderJob Creates a reminder job record to get run at some point.
func CreateReminderJob(db *gorm.DB) error {
	runId := fmt.Sprintf("Reminder Job %s", time.Now().Local())
	if err := db.Create(&jobmine.JobRecord{
		JobType:  RemindOnboardJob,
		RunId:    runId,
		Metadata: map[string]interface{}{},
		Status:   jobmine.STATUS_CREATED,
	}).Error; err != nil {
		return err
	}
	return nil
}
