package sessions

import (
	"letstalk/server/core/errs"
	"time"

	"github.com/mijia/modelq/gmq"
)

type CompositeSessionManager struct {
	mem  *InMemorySessionManager
	dbSM *DatabaseSessionManager
}

func CreateCompositeSessionManager(db *gmq.Db) ISessionManager {
	sm := CompositeSessionManager{
		CreateInMemorySessionManager(),
		CreateDBSessionManager(db),
	}
	return sm
}

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

func (sm CompositeSessionManager) CreateNewSessionForUserId(
	userId int,
) (*SessionData, errs.Error) {
	// create in both the in memory and db session manager
	return sm.CreateNewSessionForUserIdWithExpiry(userId, defaultExpiry)
}

func (sm CompositeSessionManager) CreateNewSessionForUserIdWithExpiry(
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

func (sm CompositeSessionManager) GetSessionForSessionId(
	sessionId string,
) (*SessionData, errs.Error) {
	session, ok := sm.SessionIdMapping[sessionId]
	if ok != true {
		return nil, errs.NewClientError("No session found.")
	}
	return session, nil
}

func (sm CompositeSessionManager) GetUserSessions(
	userId int,
) ([]*SessionData, errs.Error) {
	sessions, ok := sm.UserIdToSessions[userId]
	if !ok {
		return nil, errs.NewClientError("No sessions for this userId")
	}
	return sessions, nil
}
