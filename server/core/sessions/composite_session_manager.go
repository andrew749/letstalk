package sessions

import (
	"letstalk/server/core/errs"
	"time"

	"github.com/mijia/modelq/gmq"
)

type CompositeSessionManager struct {
	sessionStores []ISessionStore
}

func CreateCompositeSessionManager(db *gmq.Db, sessionManagers ...ISessionStore) ISessionManagerBase {
	sm := CompositeSessionManager{
		make([]ISessionStore, 0),
	}
	for _, x := range sessionManagers {
		sm.sessionStores = append(sm.sessionStores, x)
	}
	return sm
}

func (sm CompositeSessionManager) forEverySm(
	fx func(x ISessionStore) error,
) *errs.CompositeError {
	var compositeError *errs.CompositeError
	for _, x := range sm.sessionStores {
		errs.AppendNullableError(compositeError, fx(x))
	}
	return compositeError
}

/**
 * Stop operating once predicate is true
 */
func (sm CompositeSessionManager) forEverySmPredicate(
	fx func(x ISessionStore) (bool, error),
) *errs.CompositeError {
	var compositeError *errs.CompositeError
	for _, x := range sm.sessionStores {
		res, err := fx(x)
		errs.AppendNullableError(compositeError, err)
		if res {
			return compositeError
		}
	}
	return compositeError
}

func (sm CompositeSessionManager) AddNewSession(session *SessionData) error {
	err := sm.forEverySm(func(x ISessionStore) error {
		return x.AddNewSession(session)
	})
	return err
}

func (sm CompositeSessionManager) CreateNewSessionForUserId(
	userId int,
) (*SessionData, errs.Error) {
	defaultExpiry := time.Now().Add(time.Duration(DEFAULT_EXPIRY) * time.Hour)
	return sm.CreateNewSessionForUserIdWithExpiry(userId, defaultExpiry)
}

func (sm CompositeSessionManager) CreateNewSessionForUserIdWithExpiry(
	userId int,
	expiry time.Time,
) (*SessionData, errs.Error) {
	session, err := CreateSessionData(userId, expiry)
	if err != nil {
		return nil, errs.NewInternalError("Unable to create new session")
	}

	// maintain mappings
	sm.AddNewSession(session)
	return session, nil
}

func (sm CompositeSessionManager) GetSessionForSessionId(
	sessionId string,
) (*SessionData, errs.Error) {
	var session *SessionData
	err := sm.forEverySmPredicate(func(x ISessionStore) (bool, error) {
		res, err := x.GetSessionForSessionId(sessionId)
		if err != nil {
			return false, err
		}
		if res != nil {
			session = res
			return true, nil
		}
		return true, nil
	})

	if err != nil {
		return nil, errs.NewClientError("Unable to get session: %s", err)
	}

	return session, nil
}

func (sm CompositeSessionManager) GetUserSessions(
	userId int,
) ([]SessionData, errs.Error) {
	var res []SessionData = make([]SessionData, 0)
	err := sm.forEverySm(func(x ISessionStore) error {
		sessions, err := x.GetUserSessions(userId)
		if err != nil {
			return err
		}
		for _, session := range sessions {
			res = append(res, session)
		}
		return nil
	})
	if err != nil {
		return nil, errs.NewClientError("Could not get sessions for user %s", err)
	}
	return res, nil
}
