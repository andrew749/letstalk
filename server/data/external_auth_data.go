package data

import (
	"github.com/jinzhu/gorm"
)

type ExternalAuthData struct {
	User          User    `gorm:"foreignkey:UserId"`
	UserId        TUserID `gorm:"primary_key;not null"`
	FbUserId      *string `gorm:"null"`
	FbProfileLink *string `gorm:"null"`
	PhoneNumber   *string `gorm:"null"`
}

func CreateExternalAuthData(
	db *gorm.DB,
	userId TUserID,
	fbUserId *string,
	fbProfileLink *string,
	phoneNumber *string,
) (*ExternalAuthData, error) {
	externalAuthRecord := ExternalAuthData{
		UserId:        userId,
		FbUserId:      fbUserId,
		FbProfileLink: fbProfileLink,
		PhoneNumber:   phoneNumber,
	}

	if err := db.Create(&externalAuthRecord).Error; err != nil {
		return nil, err
	}

	return &externalAuthRecord, nil
}
