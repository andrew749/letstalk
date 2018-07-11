package sessions

import (
	"letstalk/server/core/test"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
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

func TestDatabaseSessionStore(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				sd := CreateDBSessionStore(db)
				sm := CreateCompositeSessionManager(sd)
				session1, _ := sm.CreateNewSessionForUserId(1, nil)
				session, _ := sm.GetSessionForSessionId(*session1.SessionId)
				assert.Equal(t, session1.SessionId, session.SessionId)
			},
			TestName: "Test Database session creation",
		},
	}
	test.RunTestsWithDb(tests)
}

func TestWriteThroughCache(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				sd := CreateDBSessionStore(db)
				sessionData, err := CreateSessionData(1, nil, time.Now())
				assert.NoError(t, err)

				sd.AddNewSession(sessionData)
				ss := CreateInMemorySessionStore()
				sm := CreateCompositeSessionManager(ss, sd)
				session, err := sm.GetSessionForSessionId(*sessionData.SessionId)
				assert.NoError(t, err)
				assert.NotNil(t, session)
				session2, err := ss.GetSessionForSessionId(*sessionData.SessionId)
				assert.NoError(t, err)
				assert.NotNil(t, session2)
			},
			TestName: "Test Database session creation",
		},
	}
	test.RunTestsWithDb(tests)
}
