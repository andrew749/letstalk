package sessions

import (
	"database/sql"
	"errors"
	"letstalk/server/data"
	"time"

	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
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
		if session.NotificationToken != nil {
			rlog.Debug("Storing notification data")
			notificationModel := data.NotificationTokens{
				Id:     *session.SessionId,
				UserId: session.UserId,
				Token:  *session.NotificationToken,
			}
			_, err := notificationModel.Insert(tx)
			if err != nil {
				// log the error when storing notification token
				rlog.Error(err)
			}
		}
		return nil
	})
	return err
}

func (sm DatabaseSessionStore) GetSessionForSessionId(
	sessionId string,
) (*SessionData, error) {

	find_sessions_statement, err := sm.DB.Prepare(
		`	SELECT
				session_id, sessions.user_id, token, expiry_date
			FROM
				sessions
			LEFT JOIN
				notification_tokens ON sessions.user_id=notification_tokens.user_id
			WHERE
				sessions.session_id=?`,
	)
	if err != nil {
		return nil, err
	}
	res, err := find_sessions_statement.Query(sessionId)

	if err != nil {
		return nil, err
	}

	data, err := getSessionData(res)

	if err != nil {
		return nil, err
	}

	return data[0], nil
}

func getSessionData(res *sql.Rows) ([]*SessionData, error) {
	var (
		sessionId         string
		userId            int
		notificationToken *string
		expiryDate        time.Time
	)
	result := make([]*SessionData, 0)
	for res.Next() {
		err := res.Scan(&sessionId, &userId, &notificationToken, &expiryDate)
		if err != nil {
			return nil, err
		}
		result = append(result, &SessionData{&sessionId, userId, notificationToken, expiryDate})
	}

	return result, nil
}

func (sm DatabaseSessionStore) GetUserSessions(
	userId int,
) ([]*SessionData, error) {
	find_sessions_statement, err := sm.DB.Prepare(
		`	SELECT
				session_id, sessions.user_id, token, expiry_date
			FROM
				sessions
			LEFT JOIN
				notification_tokens ON sessions.user_id=notification_tokens.user_id
			WHERE
				sessions.user_id=?`,
	)
	if err != nil {
		return nil, errors.New("Unable to perform session search operation")
	}
	sessionData, err := find_sessions_statement.Query(userId)

	if err != nil {
		return nil, errors.New("Could not get session for user")
	}

	transformedSessions, err := getSessionData(sessionData)
	if err != nil {
		return nil, errors.New("Unable to parse result")
	}

	return transformedSessions, nil
}
