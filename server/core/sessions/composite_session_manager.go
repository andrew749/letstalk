package sessions

import (
	"errors"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type CompositeSessionManager struct {
	db            *gorm.DB
	sessionStores []ISessionStore
}

func CreateCompositeSessionManager(db *gorm.DB, sessionManagers ...ISessionStore) ISessionManagerBase {
	sm := CompositeSessionManager{
		db,
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
		rlog.Debug("Adding to session store")
		return x.AddNewSession(session)
	})
	return err
}

func (sm CompositeSessionManager) CreateNewSessionForUserId(
	userId data.TUserID,
	notificationToken *string,
) (*SessionData, error) {
	defaultExpiry := time.Now().Add(time.Duration(DEFAULT_EXPIRY) * time.Hour)
	return sm.CreateNewSessionForUserIdWithExpiry(userId, notificationToken, defaultExpiry)
}

func (sm CompositeSessionManager) CreateNewSessionForUserIdWithExpiry(
	userId data.TUserID,
	notificationToken *string,
	expiry time.Time,
) (*SessionData, error) {
	session, err := CreateSessionData(userId, expiry)
	if err != nil {
		return nil, errors.New("Unable to create new session")
	}

	// store device
	if notificationToken != nil {
		if err := data.AddExpoDeviceTokenforUser(sm.db, userId, *notificationToken); err != nil {
			return nil, errors.New("Unable to register device in db.")
		}
	}

	// maintain mappings
	sm.AddNewSession(session)
	return session, nil
}

func (sm CompositeSessionManager) GetSessionForSessionId(
	sessionId string,
) (*SessionData, error) {
	var session *SessionData
	var emptyStore = false
	err := sm.forEverySmPredicate(func(x ISessionStore) (bool, error) {
		res, err := x.GetSessionForSessionId(sessionId)
		if err != nil {
			emptyStore = true
			return false, err
		}
		if res != nil {
			session = res
			return true, nil
		}
		return true, nil
	})

	if err != nil || session == nil {
		return nil, errors.New("Unable to get session")
	}

	if emptyStore {
		// FIXME: ignore for now if there is an error
		rlog.Debug("Updating session across stores")
		_ = sm.AddNewSession(session)
	}

	return session, nil
}

func (sm CompositeSessionManager) GetUserSessions(
	userId data.TUserID,
) ([]*SessionData, error) {
	var res []*SessionData = make([]*SessionData, 0)
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
		return nil, errors.New("Could not get sessions for user")
	}
	return res, nil
}
