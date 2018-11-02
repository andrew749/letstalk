package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetSurvey gets the most up-to-date survey and responses for the auth user.
func GetSurvey(c *ctx.Context) errs.Error {
	group := c.GinContext.Param("group")
	// Fetch user's survey information
	if userSurvey, err := getSurvey(c.Db, c.SessionData.UserId, data.SurveyGroup(group)); err != nil {
		return err
	} else {
		c.Result = userSurvey
		return nil
	}
}

func getSurvey(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
) (*api.Survey, errs.Error) {
	userSurvey := survey.GetSurveyDefinitionByGroup(group)
	if userSurvey == nil {
		return nil, errs.NewNotFoundError("no survey for user group '%v'", group)
	}
	if responses, err := getSurveyResponses(db, userId, group); err != nil {
		return nil, err
	} else if responses != nil {
		userSurvey.Responses = responses
	}
	return userSurvey, nil
}

func getSurveyResponses(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
) (*data.SurveyResponses, errs.Error) {
	if userSurvey, err := query.GetUserSurvey(db, userId, group); err != nil {
		return nil, errs.NewDbError(err)
	} else if userSurvey == nil {
		return nil, nil
	} else {
		return &userSurvey.Responses, nil
	}
}

// PostSurveyResponses saves a response to an onboarding survey.
func PostSurveyResponses(c *ctx.Context) errs.Error {
	var input api.Survey
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if input.Responses == nil {
		return errs.NewRequestError("Expected non-nil survey responses")
	}
	if err := saveSurveyResponses(
		c.Db, c.SessionData.UserId,
		input.Group, input.Version,
		*input.Responses,
	); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func saveSurveyResponses(
	db *gorm.DB,
	userId data.TUserID,
	group data.SurveyGroup,
	version int,
	responses data.SurveyResponses,
) errs.Error {
	var newSurvey *data.UserSurvey
	if oldSurvey, err := query.GetUserSurvey(db, userId, group); err != nil {
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
