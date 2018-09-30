package sessions

import (
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
)

type ISessionStore interface {
	GetSessionForSessionId(sessionId string) (*SessionData, error)
	GetUserSessions(userId data.TUserID) ([]*SessionData, error)
	AddNewSession(session *SessionData) error
}

type ISessionManagerBase interface {
	CreateNewSessionForUserId(userId data.TUserID) (*SessionData, error)
	CreateNewSessionForUserIdWithExpiry(userId data.TUserID, expiry time.Time) (*SessionData, error)
	GetSessionForSessionId(sessionId string) (*SessionData, error)
	GetUserSessions(userId data.TUserID) ([]*SessionData, error)
}

func CreateSessionManager(db *gorm.DB) ISessionManagerBase {
	return CreateCompositeSessionManager(
		db,
		CreateInMemorySessionStore(),
		CreateDBSessionStore(db),
	)
}

// default expiry time in days
const DEFAULT_EXPIRY = 20 * 7 * 24
