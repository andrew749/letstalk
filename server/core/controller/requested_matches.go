package controller

import (
	"fmt"
	"math/rand"
	"time"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/data"

	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type ResolveType int

type userWithCredentialId struct {
	userId       data.TUserID
	credentialId uint
}

const (
	RESOLVE_TYPE_ASKER ResolveType = iota
	RESOLVE_TYPE_ANSWERER
)

func Shuffle(vals []userWithCredentialId) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(vals) > 0 {
		n := len(vals)
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
		vals = vals[:n-1]
	}
}

func getPotentialMatchUserIds(
	db *gorm.DB,
	userId data.TUserID,
	resolveType ResolveType,
	credentialId uint,
) ([]userWithCredentialId, errs.Error) {
	if resolveType == RESOLVE_TYPE_ASKER {
		var userRequest data.UserCredentialRequest
		err := db.Where(
			"credential_id = ? and user_id = ?", credentialId, userId,
		).First(&userRequest).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		var userCredentials []data.UserCredential
		err = db.Where(
			"credential_id = ? and user_id <> ?",
			credentialId,
			userId,
		).Find(&userCredentials).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		userCredentialIds := make([]userWithCredentialId, len(userCredentials))
		for i, userCredential := range userCredentials {
			userCredentialIds[i] = userWithCredentialId{
				userCredential.UserId,
				uint(userCredential.CredentialId),
			}
		}
		return userCredentialIds, nil
	} else if resolveType == RESOLVE_TYPE_ANSWERER {
		var userCredential data.UserCredential
		err := db.Where(
			"credential_id = ? and user_id = ?", credentialId, userId,
		).First(&userCredential).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		var userRequests []data.UserCredentialRequest
		err = db.Where(
			"credential_id = ? and user_id <> ?",
			credentialId,
			userId,
		).Find(&userRequests).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		userCredentialIds := make([]userWithCredentialId, len(userRequests))
		for i, userRequest := range userRequests {
			userCredentialIds[i] = userWithCredentialId{
				userRequest.UserId,
				uint(userRequest.CredentialId),
			}
		}
		return userCredentialIds, nil
	} else {
		return nil, errs.NewRequestError("invalid resolveType %d", resolveType)
	}
}

func sendNotifications(
	c *ctx.Context,
	askerId data.TUserID,
	answererId data.TUserID,
	credentialId uint,
	name string,
) errs.Error {
	err1 := notifications.RequestToMatchNotification(
		c.Db,
		askerId,
		notifications.REQUEST_TO_MATCH_SIDE_ASKER,
		credentialId,
		name,
	)
	err2 := notifications.RequestToMatchNotification(
		c.Db,
		answererId,
		notifications.REQUEST_TO_MATCH_SIDE_ANSWERER,
		credentialId,
		name,
	)
	if err1 != nil {
		raven.CaptureError(err1, nil)
	}
	if err2 != nil {
		raven.CaptureError(err2, nil)
	}
	return nil
}

func ResolveRequestToMatchWithDelay(
	c *ctx.Context,
	resolveType ResolveType,
	credentialId uint,
	delayMs int,
) {
	<-time.After(time.Duration(delayMs) * time.Millisecond)
	err := ResolveRequestToMatch(c, resolveType, credentialId)
	if err != nil {
		rlog.Error(err)
	}
}

// TODO: This should run in a job
func ResolveRequestToMatch(
	c *ctx.Context,
	resolveType ResolveType,
	credentialId uint,
) errs.Error {
	userId := c.SessionData.UserId

	var credential data.Credential
	if err := c.Db.Where("id = ?", credentialId).First(&credential).Error; err != nil {
		return errs.NewDbError(err)
	}

	userCredentialIds, err := getPotentialMatchUserIds(c.Db, userId, resolveType, credentialId)
	if err != nil {
		return err
	}

	if len(userCredentialIds) > 0 {
		Shuffle(userCredentialIds)

		var (
			askerId    data.TUserID
			answererId data.TUserID
		)
		if resolveType == RESOLVE_TYPE_ASKER {
			askerId = userId
			answererId = userCredentialIds[0].userId
		} else if resolveType == RESOLVE_TYPE_ANSWERER {
			askerId = userCredentialIds[0].userId
			answererId = userId
		} else {
			return errs.NewRequestError("invalid resolveType %d", resolveType)
		}

		tx := c.Db.Begin()

		rlog.Info(fmt.Sprintf("Cred: %d, User: %d", credentialId, askerId))

		dbErr := tx.Where("credential_id = ? and user_id = ?", credentialId, askerId).Delete(
			data.UserCredentialRequest{},
		).Error
		if dbErr != nil {
			tx.Rollback()
			return errs.NewDbError(dbErr)
		}

		match := data.RequestMatching{
			Asker:        askerId,
			Answerer:     answererId,
			CredentialId: credentialId,
		}

		dbErr = tx.Create(&match).Error
		if dbErr != nil {
			tx.Rollback()
			return errs.NewDbError(dbErr)
		}

		tx.Commit()

		err = sendNotifications(
			c,
			askerId,
			answererId,
			credentialId,
			credential.Name,
		)
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}
