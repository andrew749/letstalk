package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func getUserSurvey(
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

	for i, userGroupSurvey := range userGroupSurveys {
		if userSurvey, ok := userSurveysByGroup[userGroupSurvey.Survey.Group]; ok {
			relevantResponses := getRelevantResponses(userGroupSurveys[i].Survey, userSurvey.Responses)
			userGroupSurveys[i].Survey.Responses = &relevantResponses
		}
	}

	return userGroupSurveys, nil
}

func GetSurvey(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
) (*api.Survey, errs.Error) {
	theSurvey := survey.GetSurveyDefinitionByGroup(group)
	if theSurvey == nil {
		return nil, errs.NewNotFoundError("no survey for user group '%v'", group)
	}
	if responses, err := getSurveyResponses(db, userId, group); err != nil {
		return nil, err
	} else if responses != nil {
		// Filter out responses that aren't relevant to the current survey version
		relevantResponses := getRelevantResponses(*theSurvey, *responses)
		theSurvey.Responses = &relevantResponses
	}
	return theSurvey, nil
}

func getRelevantResponses(
	survey api.Survey,
	responses data.SurveyResponses,
) data.SurveyResponses {
	questionKeys := make(map[data.SurveyQuestionKey]interface{})
	for _, question := range survey.Questions {
		questionKeys[question.Key] = nil
	}
	relevantResponses := make(map[data.SurveyQuestionKey]data.SurveyOptionKey)
	for key, response := range map[data.SurveyQuestionKey]data.SurveyOptionKey(responses) {
		if _, ok := questionKeys[key]; ok {
			relevantResponses[key] = response
		}
	}
	return relevantResponses
}

func getSurveyResponses(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
) (*data.SurveyResponses, errs.Error) {
	if userSurvey, err := getUserSurvey(db, userId, group); err != nil {
		return nil, errs.NewDbError(err)
	} else if userSurvey == nil {
		return nil, nil
	} else {
		return &userSurvey.Responses, nil
	}
}

func SaveSurveyResponses(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
	version int,
	responses data.SurveyResponses,
) errs.Error {
	var newSurvey *data.UserSurvey
	if oldSurvey, err := getUserSurvey(db, userId, group); err != nil {
		return errs.NewDbError(err)
	} else if oldSurvey != nil {
		newSurvey = oldSurvey
	} else {
		newSurvey = &data.UserSurvey{}
	}
	newSurvey.UserId = userId
	newSurvey.Group = group
	newSurvey.Version = version
	newSurvey.Responses = responses
	if err := db.Save(&newSurvey).Error; err != nil {
		return errs.NewInternalError("Error saving survey responses: %v", err)
	}
	return nil
}
