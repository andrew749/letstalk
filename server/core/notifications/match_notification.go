package notifications

import (
	"fmt"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// SPECIFIC NOTIFICATION MESSAGES

type RequestToMatchSide string

const (
	REQUEST_TO_MATCH_SIDE_ASKER    RequestToMatchSide = "ASKER"
	REQUEST_TO_MATCH_SIDE_ANSWERER RequestToMatchSide = "ANSWERER"
)

func RequestToMatchNotification(
	db *gorm.DB,
	recipient data.TUserID,
	side RequestToMatchSide,
	requestId uint,
	name string,
) error {
	var (
		extraData map[string]string = map[string]string{"side": string(side), "requestId": string(requestId)}
		title     string            = "You got a match!"
		message   string            = fmt.Sprintf("You got matched for \"%s\"", name)
	)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_REQUEST_TO_MATCH,
		nil,
		extraData,
	)
}

func NewMatchNotification(db *gorm.DB, recipient data.TUserID, message string) error {
	title := "You got a match!"
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_NEW_MATCH,
		nil,
		map[string]string{},
	)
}

func NewMentorNotification(db *gorm.DB, recipient data.TUserID) error {
	return NewMatchNotification(db, recipient, "You were matched with a new mentor.")
}

func NewMenteeNotification(db *gorm.DB, recipient data.TUserID) error {
	return NewMatchNotification(db, recipient, "You were matched with a new mentee.")
}

func MatchVerifiedNotification(db *gorm.DB, recipient data.TUserID, userName string) error {
	title := "You verified a match!"
	message := fmt.Sprintf("Your match with %s is now verified.", userName)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_MATCH_VERIFIED,
		nil,
		map[string]string{},
	)
}
