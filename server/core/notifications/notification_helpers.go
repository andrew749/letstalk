package notifications

import (
	"encoding/json"
	"letstalk/server/aws_utils"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/queue/queues/notification_queue"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateAndSendNotification Creates a notification object and saves to data store. Also sends to sqs so it can be processed later
func CreateAndSendNotification(
	db *gorm.DB,
	title,
	message string,
	recipient data.TUserID,
	class data.NotifType,
	thumbnail *string,
	metadata map[string]string,
) error {
	currentTime := time.Now()
	var err error
	notification, err := CreateNotification(db, recipient, class, title, message, thumbnail, currentTime, metadata)
	if err != nil {
		return err
	}

	sqsHelper, err := aws_utils.GetSQSServiceClient()
	if err != nil {
		return err
	}

	// push to sqs
	// TODO: if doesn't send then try again or set some bit saying that it wasnt sent so we can have a job retry
	return notification_queue.PushNotificationToQueue(sqsHelper, *notification)
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
		Title:  title,
	}

	if err := db.Create(dataNotif).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	return dataNotif, nil
}
