package data

import (
	"github.com/jinzhu/gorm"
)

// NOTE(wojtechnology): Right now, this is a pretty limiting data model so I see us moving off of
// this in the near future, to have one row in a table per contact info per user (labelled by
// contact info types). This will allow users to add multiple phone numbers and also will scale
// better if we want to add more different kinds of social networks.
type ExternalAuthData struct {
	User          User    `gorm:"foreignkey:UserId"`
	UserId        TUserID `gorm:"primary_key;not null"`
	FbUserId      *string `gorm:"null;size:100"`
	FbProfileLink *string `gorm:"null;type:text"`
	PhoneNumber   *string `gorm:"null;size:100"`
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
