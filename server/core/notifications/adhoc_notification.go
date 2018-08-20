package notifications

import (
	"encoding/json"
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateAdHocNotification Creates an adhoc notification as well as a page to render when users click though on the notification.
func CreateAdHocNotification(db *gorm.DB, recipient data.TUserID, title string, message string, thumbnail *string, templatePath string, templateParams map[string]string) error {
	creationTime := time.Now()
	var err error
	tx := db.Begin()
	notification, err := CreateNotification(db, recipient, NOTIF_TYPE_ADHOC, title, message, thumbnail, creationTime, templateParams)
	if err != nil {
		tx.Rollback()
		return err
	}
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
	return tx.Commit().Error
}

func createTestNotificationPage(db *gorm.DB, userId uint) error {
	d := []byte("{\"title\":\"Title\", \"body\":\"This is a body\"}")
	n, _ := CreateNotification(db, data.TUserID(userId), NOTIF_TYPE_ADHOC, "Test Notification", "Test Message", nil, time.Now(), nil)
	req := data.NotificationPage{
		NotificationId: n.ID,
		UserId:         data.TUserID(userId),
		TemplateLink:   "sample_template.html",
		Attributes:     data.JSONBlob(d),
	}
	return db.Save(&req).Error
}
