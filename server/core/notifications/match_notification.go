package notifications

import (
	"fmt"

	"letstalk/server/core/linking"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// SPECIFIC NOTIFICATION MESSAGES

type RequestToMatchSide string

const (
	REQUEST_TO_MATCH_SIDE_ASKER    RequestToMatchSide = "ASKER"
	REQUEST_TO_MATCH_SIDE_ANSWERER RequestToMatchSide = "ANSWERER"
)

func setImageUrlIfExists(extraData map[string]string, profilePic *string) {
	const IMAGE_URL = "imageUrl"
	if profilePic != nil {
		extraData[IMAGE_URL] = *profilePic
	}
}

func RequestToMatchNotification(
	db *gorm.DB,
	recipient data.TUserID,
	side RequestToMatchSide,
	matchData *data.User,
	requestId uint,
	name string,
) error {
	var (
		extraData map[string]string = map[string]string{"side": string(side), "requestId": string(requestId)}
		title     string            = "You got a match!"
		message   string            = fmt.Sprintf("You got matched for \"%s\"", name)
	)
	setImageUrlIfExists(extraData, matchData.ProfilePic)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_REQUEST_TO_MATCH,
		nil,
		extraData,
		linking.GetMatchProfileUrl(matchData.UserId),
	)
}

func newMatchNotification(db *gorm.DB, recipient data.TUserID, matchData *data.User, title string, message string) error {
	link := linking.GetMatchProfileUrl(matchData.UserId)
	extraData := make(map[string]string)
	setImageUrlIfExists(extraData, matchData.ProfilePic)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_NEW_MATCH,
		nil,
		extraData,
		link,
	)
}

// NewMentorNotification tells a user that they have a new mentor.
func NewMentorNotification(db *gorm.DB, recipient data.TUserID, mentor *data.User) error {
	title := "You have a new mentor!"
	message := fmt.Sprintf("You've been matched with a new mentor: %s %s", mentor.FirstName, mentor.LastName)
	return newMatchNotification(db, recipient, mentor, title, message)
}

// NewMenteeNotification tells a user they have a new mentee.
func NewMenteeNotification(db *gorm.DB, recipient data.TUserID, mentee *data.User) error {
	title := "You have a new mentee!"
	message := fmt.Sprintf("You've been matched with a new mentee: %s %s", mentee.FirstName, mentee.LastName)
	return newMatchNotification(db, recipient, mentee, title, message)
}

func MatchVerifiedNotification(db *gorm.DB, recipient data.TUserID, userName string, userId data.TUserID) error {
	title := "You verified a match!"
	message := fmt.Sprintf("Your match with %s is now verified.", userName)
	link := linking.GetMatchProfileUrl(userId)
	return CreateAndSendNotification(
		db,
		title,
		message,
		recipient,
		data.NOTIF_TYPE_MATCH_VERIFIED,
		nil,
		map[string]string{},
		link,
	)
}
