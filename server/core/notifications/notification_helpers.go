package notifications

import (
	"encoding/json"
	"letstalk/server/aws_utils"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/notifications"
	"letstalk/server/queue/queues/notification_queue"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateAndSendNotification Creates a notification object and saves to data store. Also sends to sqs so it can be processed later
func CreateAndSendNotification(
	db *gorm.DB,
	message string,
	recipient data.TUserID,
	class data.NotifType,
	thumbnail *string,
	metadata map[string]string,
) error {
	currentTime := time.Now()
	var err error
	notification, err := CreateNotification(db, recipient, class, message, thumbnail, currentTime, metadata)
	if err != nil {
		return err
	}

	if err = db.Save(notification).Error; err != nil {
		return err
	}

	var sendNotification = &notifications.Notification{}
	sendNotification = sendNotification.FromNotificationDataModel(*notification)

	sqsHelper, err := aws_utils.GetSQSServiceClient()
	if err != nil {
		return err
	}

	// push to sqs
	// TODO: if doesn't send then try again or set some bit saying that it wasnt sent so we can have a job retry
	return notification_queue.PushNotificationToQueue(sqsHelper, *sendNotification)
}

func SendTestNotification(message string, recipient string) error {
	notification := notifications.Notification{
		To:    recipient,
		Title: message,
	}

	sqsHelper, err := aws_utils.GetSQSServiceClient()
	if err != nil {
		return err
	}
	return notification_queue.PushNotificationToQueue(sqsHelper, notification)
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

func CreateNotification(
	db *gorm.DB,
	userId data.TUserID,
	tpe data.NotifType,
	message string,
	thumbnail *string,
	createdAt time.Time,
	dataMap map[string]string,
) (*data.Notification, errs.Error) {
	notifData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	dataNotif := &data.Notification{
		UserId: userId,
		Type:   tpe,
		State:  data.NOTIF_STATE_UNREAD,
		Data:   notifData,
	}

	if err := db.Create(dataNotif).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	// NOTE: This should not error since we just marshalled the data, so didn't add any logic for
	// deleting corrupt notifications.
	// apiNotif, err := query.NotificationDataToApi(*dataNotif)
	// if err != nil {
	// 	return nil, errs.NewInternalError(err.Error())
	// }

	// return apiNotif, nil
	return dataNotif, nil
}
