package customized_notifications

import (
	"letstalk/server/core/notifications"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

const (
	SIGNUP_NOTIFICATION_TITLE    = "Welcome to Hive!"
	SIGNUP_NOTIFICATION_MESSAGE  = "Get ready to join a new sort of network."
	SIGNUP_NOTIFICATION_TEMPLATE = "signup_notification.html"
)

var (
	SIGNUP_RUN_ID = "Signup"
)

// SendSignupNotifiction Send a notification on signup
func SendSignupNotifiction(db *gorm.DB, userId data.TUserID) error {
	return notifications.CreateAdHocNotification(
		db,
		userId,
		SIGNUP_NOTIFICATION_TITLE,
		SIGNUP_NOTIFICATION_MESSAGE,
		nil,
		SIGNUP_NOTIFICATION_TEMPLATE,
		map[string]interface{}{},
		&SIGNUP_RUN_ID,
	)
}
