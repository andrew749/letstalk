package notifications

import (
	"encoding/json"
	"letstalk/server/core/converters"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/notifications"
	"letstalk/server/queue/queues/notification_queue"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// CreateAndSendNotification Creates a notification object and saves to data store. Also sends to sqs so it can be processed later
func CreateAndSendNotification(
	db *gorm.DB,
	title,
	message string,
	recipient data.TUserID,
	class data.NotifType,
	thumbnail *string,
	metadata map[string]interface{},
	link string,
	runId *string, // optional run id to identify batch where notification was sent
) error {
	currentTime := time.Now()
	var err error
	notification, err := CreateNotification(db, recipient, class, title, message, thumbnail, currentTime, metadata, link, nil)
	if err != nil {
		return err
	}

	// push to sqs
	// TODO: if doesn't send then try again or set some bit saying that it wasnt sent so we can have a job retry
	return notification_queue.PushNotificationToQueue(*notification)
}

func UpdateNotificationState(
	db *gorm.DB,
	userId data.TUserID,
	notificationIds []uint,
	state data.NotifState,
) errs.Error {
	err := db.Model(&data.Notification{}).Where("id in (?) and user_id = ?",
		notificationIds,
		userId,
	).Updates(&data.Notification{State: state}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

// CreateNotification Creates the data model for a notification
func CreateNotification(
	db *gorm.DB,
	userId data.TUserID,
	tpe data.NotifType,
	title string,
	message string,
	thumbnail *string,
	createdAt time.Time,
	dataMap map[string]interface{},
	link string,
	runId *string,
) (*data.Notification, errs.Error) {
	notifData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	dataNotif := &data.Notification{
		UserId:        userId,
		Type:          tpe,
		State:         data.NOTIF_STATE_UNREAD,
		Data:          notifData,
		Title:         title,
		Message:       message,
		ThumbnailLink: thumbnail,
		Timestamp:     &createdAt,
		Link:          &link,
		RunId:         runId,
	}

	if err := db.Create(dataNotif).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	return dataNotif, nil
}

// FromNotificationDataModel Convert a notification data model to a version that the expo API expects
func NotificationsFromNotificationDataModel(db *gorm.DB, orig data.Notification) (*[]notifications.ExpoNotification, error) {
	// create a bunch of notifications to send based on how many registered device ids the user has
	deviceIds, err := data.GetDeviceNotificationTokensForUser(db, orig.UserId)
	if err != nil {
		return nil, err
	}

	// allocate storage
	res := make([]notifications.ExpoNotification, len(*deviceIds))

	notification, err := converters.NotificationDataToApi(orig)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to convert data.Notification to api.Notification")
	}

	// create new notification for each device id
	for i, deviceId := range *deviceIds {
		res[i] = notifications.ExpoNotification{
			To:    deviceId,
			Title: orig.Title,
			Body:  orig.Message,
			Data:  notification,
		}
	}
	return &res, nil
}
