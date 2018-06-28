package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetAllCredentials(db *gorm.DB) ([]api.Credential, errs.Error) {
	var rawCredentials []data.Credential
	if err := db.Find(&rawCredentials).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentials := make([]api.Credential, len(rawCredentials))
	for i, rawCredential := range rawCredentials {
		credentials[i] = api.Credential{rawCredential.ID, rawCredential.Name}
	}

	return credentials, nil
}

func GetUserCredentialRequests(db *gorm.DB, userId int) ([]api.Credential, errs.Error) {
	var userRequests []data.UserCredentialRequest
	if err := db.Where(
		&data.UserCredentialRequest{UserId: userId},
	).Preload("Credential").Find(&userRequests).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentials := make([]api.Credential, len(userRequests))
	for i, userRequest := range userRequests {
		credentials[i] = api.Credential{
			Id:   userRequest.CredentialId,
			Name: userRequest.Credential.Name,
		}
	}

	return credentials, nil
}

func AddUserCredentialRequest(db *gorm.DB, userId int, credentialId uint) errs.Error {
	var userRequest data.UserCredentialRequest
	err := db.Where(
		&data.UserCredentialRequest{UserId: userId, CredentialId: credentialId},
	).First(&userRequest).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return errs.NewDbError(err)
	} else if err == nil {
		return errs.NewRequestError("You have already requested this credential")
	}

	userRequest = data.UserCredentialRequest{UserId: userId, CredentialId: credentialId}

	if err := db.Save(&userRequest).Error; err != nil {
		return errs.NewDbError(err)
	}

	return nil
}

func RemoveUserCredentialRequest(db *gorm.DB, userId int, credentialId uint) errs.Error {
	err := db.Where(
		&data.UserCredentialRequest{UserId: userId, CredentialId: credentialId},
	).Delete(&data.UserCredentialRequest{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

func GetUserCredentials(db *gorm.DB, userId int) ([]api.Credential, errs.Error) {
	var userCredentials []data.UserCredential
	if err := db.Where(
		&data.UserCredential{UserId: userId},
	).Preload("Credential").Find(&userCredentials).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	credentials := make([]api.Credential, len(userCredentials))
	for i, userCredential := range userCredentials {
		credentials[i] = api.Credential{
			Id:   userCredential.CredentialId,
			Name: userCredential.Credential.Name,
		}
	}

	return credentials, nil
}

func AddUserCredential(db *gorm.DB, userId int, name string) (*uint, errs.Error) {
	var credential data.Credential

	err := db.Where(&data.Credential{Name: name}).First(&credential).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			credential = data.Credential{Name: name}

			// Add credential if it doesn't already exist.
			if err := db.Save(&credential).Error; err != nil {
				return nil, errs.NewDbError(err)
			}
		} else {
			return nil, errs.NewDbError(err)
		}
	}

	var userCredential data.UserCredential
	err = db.Where(
		&data.UserCredential{UserId: userId, CredentialId: credential.ID},
	).First(&userCredential).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errs.NewDbError(err)
	} else if err == nil {
		return nil, errs.NewRequestError("You already have this credential")
	}

	userCredential = data.UserCredential{UserId: userId, CredentialId: credential.ID}
	if err := db.Save(&userCredential).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	return &credential.ID, nil
}

func RemoveUserCredential(db *gorm.DB, userId int, credentialId uint) errs.Error {
	err := db.Where(
		&data.UserCredential{UserId: userId, CredentialId: credentialId},
	).Delete(&data.UserCredential{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
