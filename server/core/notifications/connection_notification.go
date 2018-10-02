package notifications

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"letstalk/server/core/linking"
	"letstalk/server/data"
)

func ConnectionRequestedNotification(
	db *gorm.DB,
	recipient data.TUserID,
	fromUserId data.TUserID,
	fromName string,
) error {
	var (
		title   string = "New connection request"
		message string = fmt.Sprintf("You got a connection request from %s", fromName)
	)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_REQUESTED,
		nil,
		map[string]string{},
		linking.GetMatchProfileWithButtonUrl(fromUserId),
	)
}

func ConnectionAcceptedNotification(
	db *gorm.DB,
	recipient data.TUserID,
	fromUserId data.TUserID,
	fromName string,
) error {
	var (
		title   string = "Connection request accepted"
		message string = fmt.Sprintf("%s accepted your connection request", fromName)
	)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_ACCEPTED,
		nil,
		map[string]string{},
		linking.GetMatchProfileUrl(fromUserId),
	)
}
