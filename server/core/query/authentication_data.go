package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetHashForUser(db *gorm.DB, userId data.TUserID) (*string, error) {
	var authData data.AuthenticationData
	if err := db.Where(&data.AuthenticationData{UserId: userId}).First(&authData).Error; err != nil {
		return nil, err
	}
	return &authData.PasswordHash, nil
}
