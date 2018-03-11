package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetSession(t *testing.T) {
	ss := CreateInMemorySessionStore()
	sm := CreateCompositeSessionManager(ss)

	session, _ := sm.CreateNewSessionForUserId(1, nil)

	sessions, _ := sm.GetUserSessions(1)

	assert.Equal(t, session, sessions[0])
}

func TestCreateAndGetMultipleSessions(t *testing.T) {
	ss := CreateInMemorySessionStore()
	sm := CreateCompositeSessionManager(ss)

	session1, _ := sm.CreateNewSessionForUserId(1, nil)
	session2, _ := sm.CreateNewSessionForUserId(1, nil)

	sessions, _ := sm.GetUserSessions(1)

	assert.Equal(t, session1, sessions[0])
	assert.Equal(t, session2, sessions[1])
}

func TestCreateAndGetSessionBySessionId(t *testing.T) {
	ss := CreateInMemorySessionStore()
	sm := CreateCompositeSessionManager(ss)

	session1, _ := sm.CreateNewSessionForUserId(1, nil)

	session, _ := sm.GetSessionForSessionId(*session1.SessionId)
	assert.Equal(t, session1, session)
}
