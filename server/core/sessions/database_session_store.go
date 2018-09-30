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
		tx.Rollback()
		return e
	}

	if e := tx.FirstOrCreate(&sessionModel).Error; e != nil {
		tx.Rollback()
		rlog.Error(e)
		return e
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
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

	return &SessionData{
		SessionId:  &session.SessionId,
		UserId:     session.UserId,
		ExpiryDate: session.ExpiryDate,
	}, nil
}

func (sm DatabaseSessionStore) GetUserSessions(
	userId data.TUserID,
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
		transformedSessions = append(
			transformedSessions,
			&SessionData{
				&session.SessionId,
				session.UserId,
				session.ExpiryDate,
			},
		)
	}

	return transformedSessions, nil
}
