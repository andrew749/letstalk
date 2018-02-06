package sessions

import (
	"letstalk/server/core/errs"
)

type ISessionManager interface {
	GetSessionForSessionId(sessionId string) (*SessionData, errs.Error)
	GetUserSessions(userId int) ([]*SessionData, errs.Error)
	CreateNewSessionForUserId(userId int) (*SessionData, errs.Error)
}

type InMemorySessionManager struct {
	// map session Id to a particular session
	SessionIdMapping map[string]*SessionData
	// map user id to all the sessions for this user
	UserIdToSessions map[int][]*SessionData
}

func CreateSessionManager() ISessionManager {
	sm := InMemorySessionManager{
		make(map[string]*SessionData),
		make(map[int][]*SessionData),
	}
	return sm
}

func (sm InMemorySessionManager) CreateNewSessionForUserId(
	userId int,
) (*SessionData, errs.Error) {
	session, err := CreateSessionData(userId)
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

//TODO(acod): create redis backed session manager
