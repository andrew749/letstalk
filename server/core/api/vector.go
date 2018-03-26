package api

import (
	"github.com/jinzhu/gorm"
	"letstalk/server/data"
)

type UserVectors struct {
	Me  *data.UserVector `json:"me"`
	You *data.UserVector `json:"you"`
}

func GetUserVectorsById(db *gorm.DB, userId int) (*UserVectors, error) {
	userVectors := make([]data.UserVector, 0)
	if err := db.Where("user_id = ?", userId).Find(&userVectors).Error; err != nil {
		return nil, err
	}

	var (
		me  *data.UserVector
		you *data.UserVector
	)

	// gets the last me and you vectors from list
	for _, vector := range userVectors {
		preferenceType := UserVectorPreferenceType(vector.PreferenceType)
		if preferenceType == MePreference {
			me = &vector
		} else if preferenceType == YouPreference {
			you = &vector
		}
	}

	return &UserVectors{me, you}, nil
}
