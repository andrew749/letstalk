package query

import (
	"errors"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserById(db *gorm.DB, userId int) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{UserId: userId}).First(&user).RecordNotFound() {
		return nil, errors.New("Unable to find user")
	}
	return &user, nil
}

func GetUserBySecret(db *gorm.DB, secret string) (*data.User, error) {
	var user data.User
	if db.Where(&data.User{Secret: secret}).First(&user).RecordNotFound() {
		return nil, errors.New("Unable to find user")
	}
	return &user, nil
}

func GetUserProfileById(db *gorm.DB, userId int) (*data.User, error) {
	var user data.User

	if db.Where(
		&data.User{UserId: userId},
	).Preload("ExternalAuthData").Preload("AdditionalData").First(&user).RecordNotFound() {
		return nil, errors.New("Unable to find user")
	}

	return &user, nil
}
