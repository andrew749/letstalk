package sessions

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ISessionStore interface {
	GetSessionForSessionId(sessionId string) (*SessionData, error)
	GetUserSessions(userId int) ([]*SessionData, error)
	AddNewSession(session *SessionData) error
}

type ISessionManagerBase interface {
	CreateNewSessionForUserId(userId int, notificationToken *string) (*SessionData, error)
	CreateNewSessionForUserIdWithExpiry(userId int, notificationToken *string, expiry time.Time) (*SessionData, error)
	GetSessionForSessionId(sessionId string) (*SessionData, error)
	GetUserSessions(userId int) ([]*SessionData, error)
}

func CreateSessionManager(db *gorm.DB) ISessionManagerBase {
	return CreateCompositeSessionManager(
		CreateInMemorySessionStore(),
		CreateDBSessionStore(db),
	)
}

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

//TODO(acod): create redis backed session manager
//TODO(acod): create backend job to delete stale sessions
