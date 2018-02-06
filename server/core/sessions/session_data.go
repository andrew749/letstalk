package sessions

import (
	"letstalk/server/core/utility"
	"time"
)

/**
 * Stores data related to a certain session.
 */
type SessionData struct {
	SessionId  *string
	UserId     int
	ExpiryDate time.Time
}

func CreateSessionData(userId int, expiry time.Time) (*SessionData, error) {
	sessionId, err := utility.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	return &SessionData{&sessionId, userId, expiry}, nil
}
