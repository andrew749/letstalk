package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserSurvey(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
) (*data.UserSurvey, error) {
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

func GetUserGroupSurveys(db *gorm.DB, userId data.TUserID) ([]api.UserGroupSurvey, errs.Error) {
	userGroups, err := GetUserGroups(db, userId)
	if err != nil {
		return nil, err
	}
	userGroupSurveys := make([]api.UserGroupSurvey, 0)
	for _, userGroup := range userGroups {
		survey := survey.GetSurveyDefinitionByGroupId(userGroup.GroupId)
		if survey != nil {
			userGroupSurveys = append(userGroupSurveys, api.UserGroupSurvey{
				UserGroup: api.UserGroup{
					Id:        userGroup.Id,
					GroupId:   userGroup.GroupId,
					GroupName: userGroup.GroupName,
				},
				Survey: *survey,
			})
		}
	}

	var userSurveys []data.UserSurvey
	dbErr := db.Where(&data.UserSurvey{UserId: userId}).Find(&userSurveys).Error
	if dbErr != nil {
		return nil, errs.NewDbError(dbErr)
	}
	userSurveysByGroup := make(map[data.SurveyGroup]data.UserSurvey)
	for _, userSurvey := range userSurveys {
		userSurveysByGroup[userSurvey.Group] = userSurvey
	}

	for _, userGroupSurvey := range userGroupSurveys {
		if userSurvey, ok := userSurveysByGroup[userGroupSurvey.Survey.Group]; ok {
			userGroupSurvey.Survey.Responses = &userSurvey.Responses
		}
	}

	return userGroupSurveys, nil
}
