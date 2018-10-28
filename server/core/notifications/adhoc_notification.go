package notifications

import (
	"encoding/json"
	"letstalk/server/core/linking"
	"letstalk/server/data"
	"letstalk/server/queue/queues/notification_queue"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateAdHocNotification Creates an adhoc notification as well as a
// page to render when users click though on the notification.
func CreateAdHocNotification(db *gorm.DB, recipient data.TUserID, title string, message string, thumbnail *string, templatePath string, templateParams map[string]interface{}, runId *string) error {
	creationTime := time.Now()
	var err error
	tx := db.Begin()
	// note that this cant use the helper create and send notification since we
	// probably want all data written to our db before being sent to aws

	// HACK:
	// TODO: fix this hack later since we want to keep a consistent interface
	// to create a notification but don't know the notification id until
	// the notification is saved to db
	notification, err := CreateNotification(tx, recipient, data.NOTIF_TYPE_ADHOC, title, message, thumbnail, creationTime, templateParams, linking.GetNotificationViewUrl(), runId)
	if err != nil {
		tx.Rollback()
		return err
	}
	link := linking.GetAdhocLink(notification.ID)
	notification.Link = &link
	if err = tx.Save(notification).Error; err != nil {
		tx.Rollback()
		return err
	}
	// ENDHACK:

	dataString, err := json.Marshal(templateParams)
	if err != nil {
		tx.Rollback()
		return err
	}
	page := data.NotificationPage{
		NotificationId: notification.ID,
		UserId:         recipient,
		TemplateLink:   templatePath,
		Attributes:     data.JSONBlob(dataString),
	}

	if err = tx.Save(&page).Error; err != nil {
		tx.Rollback()
		return err
	}

	// push to sqs
	err = notification_queue.PushNotificationToQueue(*notification)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
