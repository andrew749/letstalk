package sessions

import (
	"letstalk/server/core/utility"
	"letstalk/server/data"
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

func SessionDataFromDBSessionData(session data.Sessions) *SessionData {
	return &SessionData{
		SessionId:  &session.SessionId,
		UserId:     session.UserId,
		ExpiryDate: session.ExpiryDate,
	}
}

func MapSessionDataFromDBSessionData(sessions []data.Sessions) []*SessionData {
	res := make([]*SessionData, 0)
	for _, session := range sessions {
		res = append(res, SessionDataFromDBSessionData(session))
	}
	return res
}
