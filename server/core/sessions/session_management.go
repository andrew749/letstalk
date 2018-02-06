package sessions

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type ISessionManager interface {
	GetSessionForSessionId(sessionId string) (*ctx.SessionData, errs.Error)
	CreateNewSessionForUserId(userId string) (*ctx.SessionData, errs.Error)
}

type InMemorySessionManager struct {
	Sessions map[string]ctx.SessionData
}

func CreateSessionManager() ISessionManager {
	sm := InMemorySessionManager{make(map[string]ctx.SessionData)}
	return sm
}

func (sm InMemorySessionManager) CreateNewSessionForUserId(
	userId string,
) (*ctx.SessionData, errs.Error) {
	session, err := ctx.CreateSessionData(userId)
	if err != nil {
		return nil, errs.NewInternalError("Unable to create new session")
	}
	sm.Sessions[*session.SessionId] = *session
	return session, nil
}

func (sm InMemorySessionManager) GetSessionForSessionId(userId string) (*ctx.SessionData, errs.Error) {
	session, ok := sm.Sessions[userId]
	if ok != true {
		return nil, errs.NewClientError("No session found.")
	}
	return &session, nil
}

var manager = CreateSessionManager()

func GetSessionManager() ISessionManager {
	return manager
}

//TODO(acod): create redis backed session manager
