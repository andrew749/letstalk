package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserSurvey(db *gorm.DB, userId data.TUserID, group data.SurveyGroup) (*data.UserSurvey, error) {
	var survey data.UserSurvey
	result := db.Where(&data.UserSurvey{UserId: userId, Group: group}).First(&survey)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &survey, nil
}
