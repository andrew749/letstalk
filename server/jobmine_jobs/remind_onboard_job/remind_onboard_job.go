package remind_onboard_job

import (
	"encoding/json"
	"io/ioutil"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/jobmine"

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

// map to the path of a data to use for each type of generic notification.
var notificationDefinitionMapping = map[ReminderType]string{
	REMINDER_TYPE_TRAIT:    "trait_notification.json",
	REMINDER_TYPE_POSITION: "position_notification.json",
	REMINDER_TYPE_BIO:      "bio_notification.json",
	REMINDER_TYPE_GROUP:    "group_notification.json",
}

var notificationDefinitions map[ReminderType]quoteNotificationSpec
var definitionsLoaded bool = false

// load notification definition and marshall into spec
func loadQuoteNotificationDefinition(notificationDefinitionPath string) (*quoteNotificationSpec, error) {
	var notificationSpec quoteNotificationSpec

	// open file for reading
	notificationSpecData, err := ioutil.ReadFile(notificationDefinitionPath)
	if err != nil {
		rlog.Errorf("Unable to read notification definition: %+v", err)
		return nil, err
	}

	// put data into struct
	if err := json.Unmarshal(notificationSpecData, &notificationSpec); err != nil {
		rlog.Errorf("Unable to unmarshall notification definition: %+v", err)
		return nil, err
	}

	return &notificationSpec, nil
}

// helper to load all notification definitions into a map so that this is cached
func loadDefinitions() {
	if definitionsLoaded {
		return
	}

	for reminderType, path := range notificationDefinitionMapping {
		definition, err := loadQuoteNotificationDefinition(path)
		if err != nil {
			panic(err)
		}
		notificationDefinitions[reminderType] = *definition
	}
	definitionsLoaded = true
}

type quoteNotificationSpec struct {
	Title        string      `json:"title" binding:"required"`
	Message      string      `json:"message" binding:"required"`
	Body         string      `json:"body" binding:"required"`
	Link         string      `json:"link" binding:"required"`
	CallToAction string      `json:"cta" binding:"required"`
	Quotes       []quoteSpec `json:"quotes" binding:"required"`
}

type quoteSpec struct {
	Body   string `json:"body" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// ReminderJobSpec The actual reminder job that defines the operations to perform
// when scheduled
var ReminderJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType: RemindOnboardJob,
	TaskSpec: jobmine.TaskSpec{
		Execute: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord) (interface{}, error) {
			rlog.Infof("Running Task with id %s", taskRecord.ID)

			// send messages
			userId, notificationType := parseTaskRecordMetadata(taskRecord)
			var (
				templatePath string = "notification_with_quote.html"
				title        string
				message      string
			)

			// get generic message for each notification.
			title = notificationDefinitions[REMINDER_TYPE_TRAIT].Title
			message = notificationDefinitions[REMINDER_TYPE_TRAIT].Message
			// special logic for each type of notification
			switch notificationType {
			case REMINDER_TYPE_TRAIT:
				break
			case REMINDER_TYPE_BIO:
				break
			case REMINDER_TYPE_POSITION:
				break
			case REMINDER_TYPE_GROUP:
				break
			}

			// Actually create the notification to send
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
		metadata := createReminderNotificationPayloadMetadataForUsers(
			*usersWithoutTraits,
			seenUsersSet,
			REMINDER_TYPE_TRAIT,
		)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutTraits, seenUsers, seenUsersSet)

		// find all users who haven't filled in bio
		usersWithoutBio, err := usersWithoutBio(db)
		if err != nil {
			rlog.Errorf("Unable to get users without bio: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without bio %v", len(*usersWithoutBio), *usersWithoutBio)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(
				*usersWithoutBio,
				seenUsersSet,
				REMINDER_TYPE_BIO,
			)...,
		)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutBio, seenUsers, seenUsersSet)

		// find all users who havent put in a position
		usersWithoutPosition, err := usersWithoutPosition(db)
		if err != nil {
			rlog.Errorf("Unable to get users without position: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without position %v", len(*usersWithoutPosition), *usersWithoutPosition)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(
				*usersWithoutPosition,
				seenUsersSet,
				REMINDER_TYPE_POSITION,
			)...,
		)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutPosition, seenUsers, seenUsersSet)

		// find all users who haven't put in group
		usersWithoutGroup, err := usersWithoutGroup(db)
		if err != nil {
			rlog.Errorf("Unable to get users without group: %+v", err)
			return nil, err
		}
		rlog.Debugf("Got %d users without group %v", len(*usersWithoutGroup), *usersWithoutGroup)
		metadata = append(
			metadata,
			createReminderNotificationPayloadMetadataForUsers(
				*usersWithoutGroup,
				seenUsersSet,
				REMINDER_TYPE_GROUP,
			)...,
		)
		seenUsers, seenUsersSet = mergeMarkSeen(*usersWithoutGroup, seenUsers, seenUsersSet)

		// convert map to slice
		res := make([]jobmine.Metadata, 0, len(metadata))
		for _, metadata := range metadata {
			res = append(res, metadata)
		}
		rlog.Debugf("Got %d total updates: %+v", len(res), res)
		return &res, nil
	},
}

// add newly added users to list and map of seen users
// essentially upsert a map and array
func mergeMarkSeen(users []data.TUserID, seenList []data.TUserID, seenSet map[data.TUserID]bool) ([]data.TUserID, map[data.TUserID]bool) {
	for _, user := range users {
		// if this user has yet to be seen add them to the set and list
		if _, ok := seenSet[user]; !ok {
			seenSet[user] = true
			seenList = append(seenList, user)
		} else {
			rlog.Debugf("User %d already seen. Not creating notification", user)
		}
	}
	return seenList, seenSet
}

func createReminderNotificationPayloadMetadataForUsers(
	users []data.TUserID,
	seenUsersSet map[data.TUserID]bool,
	reminderType ReminderType,
) []jobmine.Metadata {
	res := make([]jobmine.Metadata, 0, len(users))
	for _, user := range users {
		if _, ok := seenUsersSet[user]; ok {
			continue
		}
		res = append(res, jobmine.Metadata(packageTaskRecordMetadata(user, reminderType)))
	}
	return res
}

func usersWithoutTraits(db *gorm.DB) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Joins("left join user_simple_traits as traits on traits.user_id = users.user_id").
		Where("traits.id is NULL").
		Pluck("users.user_id", &temp).
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutBio(db *gorm.DB) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Joins("left join user_additional_data as additional_data on additional_data.user_id = users.user_id").
		Where("bio is NULL").
		Pluck("users.user_id", &temp).
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutPosition(db *gorm.DB) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Joins("left join user_positions as positions on positions.user_id = users.user_id").
		Where("positions.id is NULL").
		Pluck("users.user_id", &temp).
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

func usersWithoutGroup(db *gorm.DB) (*[]data.TUserID, error) {
	var temp []data.TUserID
	if err := db.
		Model(&data.User{}).
		Joins("left join user_groups as groups on groups.user_id = users.user_id").
		Where("groups.id is NULL").
		Pluck("users.user_id", &temp).
		Error; err != nil {
		return nil, err
	}
	return &temp, nil
}

// CreateReminderJob Creates a reminder job record to get run at some point.
func CreateReminderJob(db *gorm.DB, runId string) error {
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
