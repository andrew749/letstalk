package data

import (
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"

	"github.com/jinzhu/gorm"
)

type AuthenticationData struct {
	UserId       TUserID `json:"user_id" gorm:"not null;primary_key;"`
	User         User    `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	PasswordHash string  `json:"password_hash" gorm:"not null;type:varchar(128);"`
}

func CreateAuthenticationData(db *gorm.DB, userId TUserID, password string) (*AuthenticationData, error) {

	hashedPassword, err := utility.HashPassword(password)

	if err != nil {
		return nil, errs.NewInternalError("Unable to hash password %s", err.Error())
	}

	authData := AuthenticationData{
		UserId:       userId,
		PasswordHash: hashedPassword,
	}
	if err := db.Create(&authData).Error; err != nil {
		return nil, err
	}
	return &authData, nil
}
