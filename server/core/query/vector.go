package query

import (
	"github.com/jinzhu/gorm"
	"letstalk/server/data"
	"letstalk/server/core/api"
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
		preferenceType := api.UserVectorPreferenceType(vector.PreferenceType)
		if preferenceType == api.PREFERENCE_TYPE_ME {
			me = &vector
		} else if preferenceType == api.PREFERENCE_TYPE_YOU {
			you = &vector
		}
	}

	return &UserVectors{me, you}, nil
}
