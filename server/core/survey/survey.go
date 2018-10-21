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
 * PostSurveyResponses saves a response to an onboarding survey.
 */
func PostSurveyResponses(c *ctx.Context) errs.Error {
	var input api.Survey
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if input.Responses == nil {
		return errs.NewRequestError("Expected non-nil survey responses")
	}
	if err := saveSurveyResponses(c.Db, c.SessionData.UserId, input.Group, input.Version, *input.Responses); err != nil {
		return err
	}
	c.Result = "Ok"
	return nil
}

func saveSurveyResponses(db *gorm.DB, userId data.TUserID, group data.SurveyGroup, version int, responses data.SurveyResponses) errs.Error {
	newSurvey := data.UserSurvey{
		UserId: userId,
		Group: group,
		Version: version,
		Responses: responses,
	}
	if oldSurvey, err := query.GetUserSurvey(db, userId, group); err != nil {
		return errs.NewDbError(err)
	} else if oldSurvey != nil {
		newSurvey.ID = oldSurvey.ID
	}
	if err := db.Save(&newSurvey).Error; err != nil {
		return errs.NewInternalError("Error saving survey responses: %v", err)
	}
	return nil
}

func GetSurveyResponses(db *gorm.DB, userId data.TUserID, group data.SurveyGroup) (*data.SurveyResponses, errs.Error) {
	if userSurvey, err := query.GetUserSurvey(db, userId, group); err != nil {
		return nil, errs.NewDbError(err)
	} else if userSurvey == nil {
		return nil, nil
	} else {
		return &userSurvey.Responses, nil
	}
}
