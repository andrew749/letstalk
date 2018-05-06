package query

import (
	"math/rand"
	"time"

	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

type ResolveType int

const (
	RESOLVE_TYPE_ASKER ResolveType = iota
	RESOLVE_TYPE_ANSWERER
)

func Shuffle(vals []data.UserCredential) {
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
	isAsker ResolveType,
	credentialId CredentialId,
) ([]int, errs.Error) {
	return nil, nil
}

// TODO: This should run in a job
func ResolveRequestToMatch(
	db *gorm.DB,
	userId int,
	isAsker ResolveType,
	credentialId CredentialId,
) (bool, errs.Error) {
	var (
		userRequest data.UserCredentialRequest
		err         error
	)
	err = db.Where("id = ? and user_id = ?", credentialId, userId).First(
		&userRequest).Error
	if err != nil {
		return false, errs.NewDbError(err)
	}

	var userCredentials []data.UserCredential
	err = db.Where(
		"position_id = ? and organization_id = ?",
		userRequest.PositionId,
		userRequest.OrganizationId,
	).Find(&userCredentials).Error
	if err != nil {
		return false, errs.NewDbError(err)
	}

	if len(userCredentials) > 0 {
		Shuffle(userCredentials)
		userCredential := userCredentials[0]

		tx := db.Begin()

		err = tx.Where("id = ? and user_id = ?", credentialId, userId).Delete(
			data.UserCredentialRequest{}).Error
		if err != nil {
			tx.Rollback()
			return false, errs.NewDbError(err)
		}

		err = tx.Exec(
			"INSERT INTO answerers (user_user_id, answerer_id) VALUES (?, ?)",
			userId,
			userCredential.UserId,
		).Error
		if err != nil {
			tx.Rollback()
			return false, errs.NewDbError(err)
		}

		err = tx.Exec(
			"INSERT INTO askers (user_user_id, asker_id) VALUES (?, ?)",
			userCredential.UserId,
			userId,
		).Error
		if err != nil {
			tx.Rollback()
			return false, errs.NewDbError(err)
		}

		tx.Commit()

		return true, nil
	} else {
		return false, nil
	}
}
