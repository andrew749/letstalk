package sessions

import (
	"letstalk/server/core/errs"
	"time"
)

type InMemorySessionManager struct {
	// map session Id to a particular session
	SessionIdMapping map[string]*SessionData
	// map user id to all the sessions for this user
	UserIdToSessions map[int][]*SessionData
}

func CreateInMemorySessionManager() ISessionManager {
	sm := InMemorySessionManager{
		make(map[string]*SessionData),
		make(map[int][]*SessionData),
	}
	return sm
}

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

func (sm InMemorySessionManager) CreateNewSessionForUserId(
	userId int,
) (*SessionData, errs.Error) {
	defaultExpiry := time.Now().Add(time.Duration(DEFAULT_EXPIRY) * time.Hour)
	return sm.CreateNewSessionForUserIdWithExpiry(userId, defaultExpiry)
}

func (sm InMemorySessionManager) CreateNewSessionForUserIdWithExpiry(
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

func (sm InMemorySessionManager) GetSessionForSessionId(
	sessionId string,
) (*SessionData, errs.Error) {
	session, ok := sm.SessionIdMapping[sessionId]
	if ok != true {
		return nil, errs.NewClientError("No session found.")
	}
	return session, nil
}

func (sm InMemorySessionManager) GetUserSessions(
	userId int,
) ([]*SessionData, errs.Error) {
	sessions, ok := sm.UserIdToSessions[userId]
	if !ok {
		return nil, errs.NewClientError("No sessions for this userId")
	}
	return sessions, nil
}
