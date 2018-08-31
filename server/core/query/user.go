package query

import (
	"letstalk/server/core/errs"
	"letstalk/server/data"

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
		return nil, errs.NewNotFoundError("Unable to find user")
	}
	return &user, nil
}

func GetUserBySecret(db *gorm.DB, secret string) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{Secret: secret}).First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}
	return &user, nil
}

func GetUserProfileById(db *gorm.DB, userId data.TUserID) (*data.User, error) {
	var user data.User

	if db.Where(
		&data.User{UserId: userId},
	).Preload("ExternalAuthData").Preload("AdditionalData").
		Preload("UserPositions").Preload("UserSimpleTraits").First(&user).RecordNotFound() {
		return nil, errs.NewNotFoundError("Unable to find user")
	}

	return &user, nil
}
