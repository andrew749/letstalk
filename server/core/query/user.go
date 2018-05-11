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

func GetUserByIdWithExternalAuth(db *gorm.DB, userId int) (*data.User, error) {
	var user data.User

	if db.Where(
		&data.User{UserId: userId},
	).Preload("ExternalAuthData").First(&user).RecordNotFound() {
		return nil, errors.New("Unable to find user")
	}

	return &user, nil
}
