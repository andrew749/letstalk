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

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

func CreateSessionData(userId int) (*SessionData, error) {
	sessionId, err := utility.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}

	expiryDate := time.Now().Add(time.Duration(DEFAULT_EXPIRY) * time.Hour)

	return &SessionData{&sessionId, userId, expiryDate}, nil
}
