package sessions

import (
	"time"

	"github.com/jinzhu/gorm"
	"letstalk/server/core/errs"
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

func GetDeviceTokensForUser(manager ISessionManagerBase, userId int) ([]string, errs.Error) {
	userSessions, err := manager.GetUserSessions(userId)
	if err != nil {
		return nil, errs.NewClientError(err.Error())
	}
	uniqueDeviceTokens := make(map[string]interface{})
	for _, session := range userSessions {
		if session.NotificationToken != nil {
			uniqueDeviceTokens[*session.NotificationToken] = nil
		}
	}
	deviceTokens := make([]string, 0, len(uniqueDeviceTokens))
	for token := range uniqueDeviceTokens {
		deviceTokens = append(deviceTokens, token)
	}
	return deviceTokens, nil
}

//TODO(acod): create redis backed session manager
//TODO(acod): create backend job to delete stale sessions
