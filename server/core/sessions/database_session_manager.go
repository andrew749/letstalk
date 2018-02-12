package sessions

import (
	"letstalk/server/core/errs"
	"time"

	"github.com/mijia/modelq/gmq"
)

type DatabaseSessionManager struct {
	DB *gmq.Db
}

func CreateDBSessionManager(db *gmq.Db) ISessionManager {
	sm := DatabaseSessionManager{
		DB: db,
	}
	return sm
}

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

func (sm DatabaseSessionManager) CreateNewSessionForUserId(
	userId int,
) (*SessionData, errs.Error) {
	defaultExpiry := time.Now().Add(time.Duration(DEFAULT_EXPIRY) * time.Hour)
	return sm.CreateNewSessionForUserIdWithExpiry(userId, defaultExpiry)
}

func (sm DatabaseSessionManager) CreateNewSessionForUserIdWithExpiry(
	userId int,
	expiry time.Time,
) (*SessionData, errs.Error) {
	session, err := CreateSessionData(userId, expiry)
	if err != nil {
		return nil, errs.NewInternalError("Unable to create new session")
	}

	// maintain mappings
	sm.SessionIdMapping[*session.SessionId] = session
	sm.UserIdToSessions[userId] = append(sm.UserIdToSessions[userId], session)
	return session, nil
}

func (sm DatabaseSessionManager) GetSessionForSessionId(
	sessionId string,
) (*SessionData, errs.Error) {
	session, ok := sm.SessionIdMapping[sessionId]
	if ok != true {
		return nil, errs.NewClientError("No session found.")
	}
	return session, nil
}

func (sm DatabaseSessionManager) GetUserSessions(
	userId int,
) ([]*SessionData, errs.Error) {
	sessions, ok := sm.UserIdToSessions[userId]
	if !ok {
		return nil, errs.NewClientError("No sessions for this userId")
	}
	return sessions, nil
}
