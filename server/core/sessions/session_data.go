package sessions

import (
	"letstalk/server/core/utility"
	"time"
)

/**
 * Stores data related to a certain session.
 */
type SessionData struct {
	SessionId         *string
	UserId            int
	NotificationToken *string
	ExpiryDate        time.Time
}

func CreateSessionData(
	userId int,
	notificationToken *string,
	expiry time.Time,
) (*SessionData, error) {
	sessionId, err := utility.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	return &SessionData{&sessionId, userId, notificationToken, expiry}, nil
}
