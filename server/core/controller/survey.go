package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

// GetSurvey gets the most up-to-date survey and responses for the auth user.
func GetSurvey(c *ctx.Context) errs.Error {
	group := c.GinContext.Param("group")
	// Fetch user's survey information
	if userSurvey, err := query.GetSurvey(
		c.Db,
		c.SessionData.UserId,
		data.SurveyGroup(group),
	); err != nil {
		return err
	} else {
		c.Result = userSurvey
		return nil
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
	if err := query.SaveSurveyResponses(
		c.Db, c.SessionData.UserId,
		input.Group, input.Version,
		*input.Responses,
	); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}
