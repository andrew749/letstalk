package query

import (
	"math/rand"
	"time"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/notifications"

	"github.com/jinzhu/gorm"
)

type ResolveType int

type userWithCredentialId struct {
	userId       int
	credentialId CredentialId
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
	userId int,
	resolveType ResolveType,
	credentialId CredentialId,
) ([]userWithCredentialId, errs.Error) {
	if resolveType == RESOLVE_TYPE_ASKER {
		var userRequest data.UserCredentialRequest
		err := db.Where("id = ? and user_id = ?", credentialId, userId).First(&userRequest).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		var userCredentials []data.UserCredential
		err = db.Where(
			"position_id = ? and organization_id = ?",
			userRequest.PositionId,
			userRequest.OrganizationId,
		).Find(&userCredentials).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		userCredentialIds := make([]userWithCredentialId, len(userCredentials))
		for i, userCredential := range userCredentials {
			userCredentialIds[i] = userWithCredentialId{
				userCredential.UserId,
				CredentialId(userCredential.ID),
			}
		}
		return userCredentialIds, nil
	} else if resolveType == RESOLVE_TYPE_ANSWERER {
		var userCredential data.UserCredential
		err := db.Where("id = ? and user_id = ?", credentialId, userId).First(&userCredential).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		var userRequests []data.UserCredentialRequest
		err = db.Where(
			"position_id = ? and organization_id = ?",
			userCredential.PositionId,
			userCredential.OrganizationId,
		).Find(&userRequests).Error
		if err != nil {
			return nil, errs.NewDbError(err)
		}
		userCredentialIds := make([]userWithCredentialId, len(userRequests))
		for i, userRequest := range userRequests {
			userCredentialIds[i] = userWithCredentialId{
				userRequest.UserId,
				CredentialId(userRequest.ID),
			}
		}
		return userCredentialIds, nil
	} else {
		return nil, errs.NewClientError("invalid resolveType %d", resolveType)
	}
}

func getDeviceTokensForUser(c *ctx.Context, userId int) ([]string, errs.Error) {
	sessions, err := (*c.SessionManager).GetUserSessions(userId)
	if err != nil {
		return nil, errs.NewClientError(err.Error())
	}
	uniqueDeviceTokens := make(map[string]interface{})
	for _, session := range sessions {
		if session.NotificationToken != nil {
			uniqueDeviceTokens[*session.NotificationToken] = nil
		}
	}
	deviceTokens := make([]string, 0, len(uniqueDeviceTokens))
	for token, _ := range uniqueDeviceTokens {
		deviceTokens = append(deviceTokens, token)
	}
	return deviceTokens, nil
}

func sendNotifications(c *ctx.Context, askerId int, answererId int) errs.Error {
	askerDeviceTokens, err := getDeviceTokensForUser(c, askerId)
	if err != nil {
		return err
	}
	answererDeviceTokens, err := getDeviceTokensForUser(c, answererId)
	if err != nil {
		return err
	}
	for _, token := range askerDeviceTokens {
		notifications.CreateAndSendNotification(token, "You got a match", "Match")
	}
	for _, token := range answererDeviceTokens {
		notifications.CreateAndSendNotification(token, "You got a match", "Match")
	}
	return nil
}

// TODO: This should run in a job
func ResolveRequestToMatch(
	c *ctx.Context,
	userId int,
	resolveType ResolveType,
	credentialId CredentialId,
) errs.Error {
	userCredentialIds, err := getPotentialMatchUserIds(c.Db, userId, resolveType, credentialId)
	if err != nil {
		return err
	}

	if len(userCredentialIds) > 0 {
		Shuffle(userCredentialIds)

		var (
			askerId             int
			answererId          int
			credentialRequestId CredentialId // We only delete the credential request, not the credential
		)
		if resolveType == RESOLVE_TYPE_ASKER {
			askerId = userId
			answererId = userCredentialIds[0].userId
			credentialRequestId = credentialId
		} else if resolveType == RESOLVE_TYPE_ANSWERER {
			askerId = userCredentialIds[0].userId
			answererId = userId
			credentialRequestId = userCredentialIds[0].credentialId
		} else {
			return errs.NewClientError("invalid resolveType %d", resolveType)
		}

		tx := c.Db.Begin()

		dbErr := tx.Where("id = ? and user_id = ?", credentialRequestId, askerId).Delete(
			data.UserCredentialRequest{}).Error
		if dbErr != nil {
			tx.Rollback()
			return errs.NewDbError(dbErr)
		}

		match := data.RequestMatching{
			Asker:    askerId,
			Answerer: answererId,
		}

		dbErr = tx.Create(&match).Error
		if dbErr != nil {
			tx.Rollback()
			return errs.NewDbError(dbErr)
		}

		tx.Commit()

		err = sendNotifications(c, askerId, answererId)
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}

func GetAnswerersByAskerId(db *gorm.DB, askerId int) ([]data.RequestMatching, error) {
	matchings := make([]data.RequestMatching, 0)
	if err := db.Where(
		&data.RequestMatching{Asker: askerId},
	).Preload("AnswererUser").Preload("AnswererUser.ExternalAuthData").Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}

func GetAskersByAnswererId(db *gorm.DB, answererId int) ([]data.RequestMatching, error) {
	matchings := make([]data.RequestMatching, 0)
	if err := db.Where(
		&data.RequestMatching{Answerer: answererId},
	).Preload("AskerUser").Preload("AskerUser.ExternalAuthData").Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}
