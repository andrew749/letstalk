package query

import (
	"math/rand"
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
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

func GetAnswerersByAskerId(
	db *gorm.DB,
	askerId data.TUserID,
	flag api.ReqMatchingInfoFlag,
) ([]data.RequestMatching, error) {
	matchings := make([]data.RequestMatching, 0)
	req := db.Where(&data.RequestMatching{Asker: askerId}).Preload("AnswererUser")

	if flag&api.REQ_MATCHING_INFO_FLAG_AUTH_DATA != 0 {
		req = req.Preload("AnswererUser.ExternalAuthData")
	}
	if flag&api.REQ_MATCHING_INFO_FLAG_CREDENTIAL != 0 {
		req = req.Preload("Credential")
	}

	if err := req.Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}

func GetAskersByAnswererId(
	db *gorm.DB,
	answererId data.TUserID,
	flag api.ReqMatchingInfoFlag,
) ([]data.RequestMatching, error) {
	matchings := make([]data.RequestMatching, 0)
	req := db.Where(&data.RequestMatching{Answerer: answererId}).Preload("AskerUser")

	if flag&api.REQ_MATCHING_INFO_FLAG_AUTH_DATA != 0 {
		req = req.Preload("AskerUser.ExternalAuthData")
	}
	if flag&api.REQ_MATCHING_INFO_FLAG_CREDENTIAL != 0 {
		req = req.Preload("Credential")
	}

	if err := req.Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}

func RemoveAllMatches(db *gorm.DB, fstUserId data.TUserID, sndUserId data.TUserID) errs.Error {
	err := db.Where(
		"(answerer = ? and asker = ?) or (asker = ? and answerer = ?)",
		fstUserId,
		sndUserId,
		fstUserId,
		sndUserId,
	).Delete(data.RequestMatching{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
