package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserSurvey(db *gorm.DB, userId data.TUserID, version data.SurveyVersion) (*data.UserSurvey, error) {
	var survey data.UserSurvey
	result := db.Where(&data.UserSurvey{UserId: userId, Version: version}).First(&survey)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &survey, nil
}
