package survey

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

/**
 * PostSurveyResponses saves a set of survey responses.
 */
func PostSurveyResponses(c *ctx.Context) errs.Error {
	var input api.Survey
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if err := saveSurveyResponses(c.Db, c.SessionData.UserId, input.Version, input.Responses); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func saveSurveyResponses(db *gorm.DB, userId data.TUserID, version data.SurveyVersion, responses data.SurveyResponses) errs.Error {
	newSurvey := data.UserSurvey{}
	if oldSurvey, err := query.GetUserSurvey(db, userId, version); err != nil {
		return errs.NewDbError(err)
	} else if oldSurvey != nil {
		newSurvey.ID = oldSurvey.ID
	}
	newSurvey.UserId = userId
	newSurvey.Version = version
	newSurvey.Responses = responses
	if err := db.Save(&newSurvey).Error; err != nil {
		return errs.NewInternalError("Error saving survey responses: %v", err)
	}
	return nil
}

func GetSurveyResponses(db *gorm.DB, userId data.TUserID, version data.SurveyVersion) (*data.SurveyResponses, errs.Error) {
	if survey, err := query.GetUserSurvey(db, userId, version); err != nil {
		return nil, errs.NewDbError(err)
	} else if survey == nil {
		return nil, nil
	} else {
		return &survey.Responses, nil
	}
}
