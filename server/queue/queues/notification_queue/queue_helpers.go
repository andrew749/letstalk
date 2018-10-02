package notification_queue

import (
	"letstalk/server/constants"
	"letstalk/server/data"
	"letstalk/server/queue"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type NotificationQueueData struct {
	ID uint `json:"id"`
}

// DataNotificationModelToQueueModel Convert a data model of a notifiation to a
// serializable model that is stored in an sqs queue.
func DataNotificationModelToQueueModel(notification data.Notification) NotificationQueueData {
	return NotificationQueueData{
		ID: notification.ID,
	}
}

// QueueModelToDataNotificationModel Convert a serialized queue model to a
// data.Notification by looking up the appropriate notification in the db
func QueueModelToDataNotificationModel(db *gorm.DB, notification NotificationQueueData) (data.Notification, error) {
	var res data.Notification
	err := db.First(&res, notification.ID).Error
	return res, err
}

func PushNotificationToQueue(notification data.Notification) error {
	rlog.Debugf("%#v", notification)
	queueData := DataNotificationModelToQueueModel(notification)
	_, err := queue.AddNewMessage(utility.QueueHelper, constants.NotificationQueueID, constants.NotificationQueueUrl, queueData)
	return err
}
