package notifications

import (
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
)

func CreateAdHocNotification(db *gorm.DB, recipient data.TUserID, templatePath string) error {
	return nil
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
