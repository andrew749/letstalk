package sessions

import (
	"letstalk/server/core/errs"
	"time"
)

type ISessionManager interface {
	GetSessionForSessionId(sessionId string) (*SessionData, errs.Error)
	GetUserSessions(userId int) ([]*SessionData, errs.Error)
	CreateNewSessionForUserId(userId int) (*SessionData, errs.Error)
	CreateNewSessionForUserIdWithExpiry(userId int, expiry time.Time) (*SessionData, errs.Error)
}

//TODO(acod): create redis backed session manager
//TODO(acod): create backend job to delete stale sessions
