package sessions

import (
	"errors"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type DatabaseSessionStore struct {
	DB *gorm.DB
}

func CreateDBSessionStore(db *gorm.DB) ISessionStore {
	sm := DatabaseSessionStore{
		DB: db,
	}
	return sm
}

func (sm DatabaseSessionStore) AddNewSession(session *SessionData) error {
	sessionModel := data.Session{
		SessionId:  *session.SessionId,
		UserId:     session.UserId,
		ExpiryDate: session.ExpiryDate,
	}

	tx := sm.DB.Begin()
	if e := tx.Error; e != nil {
		return e
	}

	if e := tx.Create(sessionModel).Error; e != nil {
		return e
	}

	if session.NotificationToken != nil {
		rlog.Debug("Storing notification data")
		notificationModel := data.NotificationToken{
			SessionId: *session.SessionId,
			Token:     *session.NotificationToken,
		}

		if err := tx.FirstOrCreate(&notificationModel).Error; err != nil {
			rlog.Error(err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (sm DatabaseSessionStore) GetSessionForSessionId(
	sessionId string,
) (*SessionData, error) {

	var session data.Session
	if sm.DB.Where("session_id = ?", sessionId).First(&session).RecordNotFound() {
		return nil, errors.New("Unable to find session.")
	}
	var notificationToken *string = nil

	if session.NotificationToken != nil {
		notificationToken = &session.NotificationToken.Token
	}

	return &SessionData{
		SessionId:         &session.SessionId,
		UserId:            session.UserId,
		NotificationToken: notificationToken,
		ExpiryDate:        session.ExpiryDate,
	}, nil
}

func (sm DatabaseSessionStore) GetUserSessions(
	userId int,
) ([]*SessionData, error) {
	sessions := make([]data.Session, 0)
	if err := sm.DB.Where(
		"user_id = ?",
		userId,
	).Preload("NotificationToken").Find(&sessions).Error; err != nil {
		return nil, err
	}

	transformedSessions := make([]*SessionData, 0)
	for _, session := range sessions {
		var token *string
		if session.NotificationToken != nil {
			token = &session.NotificationToken.Token
		}
		transformedSessions = append(
			transformedSessions,
			&SessionData{
				&session.SessionId,
				session.UserId,
				token,
				session.ExpiryDate,
			},
		)
	}

	return transformedSessions, nil
}
