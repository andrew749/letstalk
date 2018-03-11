package sessions

import (
	"letstalk/server/core/errs"
	"time"

	"github.com/mijia/modelq/gmq"
)

type ISessionStore interface {
	GetSessionForSessionId(sessionId string) (*SessionData, errs.Error)
	GetUserSessions(userId int) ([]*SessionData, errs.Error)
	AddNewSession(session *SessionData) error
}

type ISessionManagerBase interface {
	CreateNewSessionForUserId(userId int, notificationToken *string) (*SessionData, errs.Error)
	CreateNewSessionForUserIdWithExpiry(userId int, notificationToken *string, expiry time.Time) (*SessionData, errs.Error)
	GetSessionForSessionId(sessionId string) (*SessionData, errs.Error)
	GetUserSessions(userId int) ([]*SessionData, errs.Error)
}

func CreateSessionManager(db *gmq.Db) ISessionManagerBase {
	return CreateCompositeSessionManager(
		CreateInMemorySessionStore(),
		CreateDBSessionStore(db),
	)
}

// default expiry time in days
const DEFAULT_EXPIRY = 7 * 24

//TODO(acod): create redis backed session manager
//TODO(acod): create backend job to delete stale sessions
