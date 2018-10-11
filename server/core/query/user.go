package query

import (
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"errors"

	"github.com/jinzhu/gorm"
)

func GetUserById(db *gorm.DB, userId data.TUserID) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{UserId: userId}).First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{Email: email}).First(&user).RecordNotFound() {
		// fallback to looking for uwaterloo email
		var verifyEmailId data.VerifyEmailId
		if db.Where(&data.VerifyEmailId{Email: email}).Preload("users").First(&verifyEmailId).RecordNotFound() {
			return nil, errs.NewNotFoundError("Unable to find user")
		}
		return &verifyEmailId.User, nil
	}
	return &user, nil
}

func GetUserBySecret(db *gorm.DB, secret string) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{Secret: secret}).First(&user).RecordNotFound() {
		return nil, errors.New("unable to find user")
	}
	return &user, nil
}

func GetUserProfileById(
	db *gorm.DB,
	userId data.TUserID,
	includeContactInfo bool,
) (*data.User, error) {
	var user data.User

	query := db.Where(
		&data.User{UserId: userId},
	).Preload("AdditionalData").Preload("UserPositions").Preload("UserSimpleTraits")

	if includeContactInfo {
		query = query.Preload("ExternalAuthData")
	}

	if query.First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}

	return &user, nil
}
