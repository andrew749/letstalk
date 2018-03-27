package api

import (
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

type UserVectors struct {
	Me  *data.UserVector `json:"me"`
	You *data.UserVector `json:"you"`
}

func GetUserVectorsById(db *gmq.Db, userId int) (*UserVectors, error) {
	vectors, err := data.UserVectorObjs.
		Select().
		Where(data.UserVectorObjs.FilterUserId("=", userId)).
		List(db)

	if err != nil {
		return nil, err
	}

	var (
		me  *data.UserVector
		you *data.UserVector
	)

	// gets the last me and you vectors from list
	for _, vector := range vectors {
		preferenceType := UserVectorPreferenceType(vector.PreferenceType)
		if preferenceType == MePreference {
			me = &vector
		} else if preferenceType == YouPreference {
			you = &vector
		}
	}

	return &UserVectors{me, you}, nil
}
