package notifications

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"

	"letstalk/server/core/linking"
	"letstalk/server/data"
)

func setConnUserId(extraData map[string]interface{}, connUserId data.TUserID) {
	extraData["connUserId"] = strconv.Itoa(int(connUserId))
}

func ConnectionRequestedNotification(
	db *gorm.DB,
	recipient data.TUserID,
	fromUserId data.TUserID,
	fromName string,
) error {
	var (
		extraData map[string]interface{} = map[string]interface{}{}
		title     string                 = "New connection request"
		message   string                 = fmt.Sprintf("You got a connection request from %s", fromName)
	)
	setConnUserId(extraData, fromUserId)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_REQUESTED,
		nil,
		extraData,
		linking.GetMatchProfileWithButtonUrl(fromUserId),
		nil,
	)
}

func ConnectionAcceptedNotification(
	db *gorm.DB,
	recipient data.TUserID,
	fromUserId data.TUserID,
	fromName string,
) error {
	var (
		extraData map[string]interface{} = map[string]interface{}{}
		title     string                 = "Connection request accepted"
		message   string                 = fmt.Sprintf("%s accepted your connection request", fromName)
	)
	setConnUserId(extraData, fromUserId)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_CONNECTION_ACCEPTED,
		nil,
		extraData,
		linking.GetMatchProfileUrl(fromUserId),
		nil,
	)
}
