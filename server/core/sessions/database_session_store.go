package sessions

import (
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

type DatabaseSessionStore struct {
	DB *gmq.Db
}

func CreateDBSessionStore(db *gmq.Db) ISessionStore {
	sm := DatabaseSessionStore{
		DB: db,
	}
	return sm
}

func (sm DatabaseSessionStore) AddNewSession(session *SessionData) error {
	sessionModel := data.Sessions{
		SessionId:  *session.SessionId,
		UserId:     session.UserId,
		ExpiryDate: session.ExpiryDate,
	}
	err := gmq.WithinTx(sm.DB, func(tx *gmq.Tx) error {
		_, e := sessionModel.Insert(tx)
		if e != nil {
			return e
		}
		return nil
	})
	return err
}

func (sm DatabaseSessionStore) GetSessionForSessionId(
	sessionId string,
) (*SessionData, errs.Error) {
	res, err := data.SessionsObjs.
		Select().
		Where(data.SessionsObjs.FilterSessionId("=", sessionId)).
		One(sm.DB)

	if err != nil {
		return nil, errs.NewClientError("Could not get session for user: %s", res)
	}
	convertedSessionData := SessionDataFromDBSessionData(res)
	return convertedSessionData, nil
}

func (sm DatabaseSessionStore) GetUserSessions(
	userId int,
) ([]*SessionData, errs.Error) {
	sessions, err := data.SessionsObjs.
		Select().
		Where(data.SessionsObjs.FilterUserId("=", userId)).
		List(sm.DB)

	if err != nil {
		return nil, errs.NewClientError("Could not get session for user: %s", sessions)
	}

	return MapSessionDataFromDBSessionData(sessions), nil
}
