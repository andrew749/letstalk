package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetSession(t *testing.T) {
	sm := CreateSessionManager()

	session, _ := sm.CreateNewSessionForUserId(1)

	sessions, _ := sm.GetUserSessions(1)

	assert.Equal(t, session, sessions[0])
}

func TestCreateAndGetMultipleSessions(t *testing.T) {
	sm := CreateSessionManager()

	session1, _ := sm.CreateNewSessionForUserId(1)
	session2, _ := sm.CreateNewSessionForUserId(1)

	sessions, _ := sm.GetUserSessions(1)

	assert.Equal(t, session1, sessions[0])
	assert.Equal(t, session2, sessions[1])
}

func TestCreateAndGetSessionBySessionId(t *testing.T) {
	sm := CreateSessionManager()

	session1, _ := sm.CreateNewSessionForUserId(1)

	session, _ := sm.GetSessionForSessionId(*session1.SessionId)
	assert.Equal(t, session1, session)
}
