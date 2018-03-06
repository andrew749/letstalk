package sessions

import (
	"database/sql"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"time"

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

	find_sessions_statement, err := sm.DB.Prepare(
		`	SELECT
				session_id, user_id, notification_token, expiry_date
			FROM
				sessions
			INNER JOIN
				notification_tokens ON sessions.user_id=notification_tokens.user_id
			WHERE
				sessions.session_id=?`,
	)
	if err != nil {
		return nil, errs.NewInternalError("Unable to perform session search operation")
	}
	res, err := find_sessions_statement.Query(sessionId)

	if err != nil {
		return nil, errs.NewClientError("Could not get session for user: %s", res)
	}

	data, err := getSessionData(res)

	if err != nil {
		return nil, errs.NewClientError("Unable to get session data")
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
) ([]*SessionData, errs.Error) {
	sessions, err := data.SessionsObjs.
		Select().
		Where(data.SessionsObjs.FilterUserId("=", userId)).
		List(sm.DB)
	find_sessions_statement, err := sm.DB.Prepare(
		`	SELECT
				session_id, user_id, notification_token, expiry_date
			FROM
				sessions
			INNER JOIN
				notification_tokens ON sessions.user_id=notification_tokens.user_id
			WHERE
				sessions.user_id=?`,
	)
	if err != nil {
		return nil, errs.NewInternalError("Unable to perform session search operation")
	}
	sessionData, err := find_sessions_statement.Query(userId)

	if err != nil {
		return nil, errs.NewClientError("Could not get session for user: %s", sessions)
	}

	transformedSessions, err := getSessionData(sessionData)
	if err != nil {
		return nil, errs.NewInternalError("Unable to parse result")
	}

	return transformedSessions, nil
}
