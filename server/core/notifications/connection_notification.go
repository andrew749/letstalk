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
	fromUserData *data.User,
	fromName string,
) error {
	var (
		extraData map[string]string = map[string]string{}
		title   string = "New connection request"
		message string = fmt.Sprintf("You got a connection request from %s", fromName)
	)
	setImageUrlIfExists(extraData, fromUserData.ProfilePic)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_REQUESTED,
		nil,
		extraData,
		linking.GetMatchProfileWithButtonUrl(fromUserData.UserId),
	)
}

func ConnectionAcceptedNotification(
	db *gorm.DB,
	recipient data.TUserID,
	fromUserData *data.User,
	fromName string,
) error {
	var (
		extraData map[string]string = map[string]string{}
		title   string = "Connection request accepted"
		message string = fmt.Sprintf("%s accepted your connection request", fromName)
	)
	setImageUrlIfExists(extraData, fromUserData.ProfilePic)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_ACCEPTED,
		nil,
		extraData,
		linking.GetMatchProfileUrl(fromUserData.UserId),
	)
}
