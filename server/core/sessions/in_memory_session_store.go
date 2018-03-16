package sessions

import (
	"errors"
)

type InMemorySessionStore struct {
	// map session Id to a particular session
	SessionIdMapping map[string]SessionData
	// map user id to all the sessions for this user
	UserIdToSessions map[int][]*SessionData
}

func CreateInMemorySessionStore() ISessionStore {
	sm := InMemorySessionStore{
		make(map[string]SessionData),
		make(map[int][]*SessionData),
	}
	return sm
}

func (sm InMemorySessionStore) AddNewSession(session *SessionData) error {
	sm.SessionIdMapping[*session.SessionId] = *session
	sm.UserIdToSessions[session.UserId] = append(
		sm.UserIdToSessions[session.UserId],
		session,
	)
	return nil
}

func (sm InMemorySessionStore) GetSessionForSessionId(
	sessionId string,
) (*SessionData, error) {
	session, ok := sm.SessionIdMapping[sessionId]
	if ok != true {
		return nil, errors.New("Unable to find session in memory")
	}
	return &session, nil
}

func (sm InMemorySessionStore) GetUserSessions(
	userId int,
) ([]*SessionData, error) {
	sessions, ok := sm.UserIdToSessions[userId]
	if !ok {
		return nil, errors.New("No sessions for this userId")
	}
	return sessions, nil
}
