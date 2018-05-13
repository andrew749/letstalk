package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func GetMyProfileController(c *ctx.Context) errs.Error {
	user, err := query.GetUserByIdWithExternalAuth(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewClientError("Unable to get user data.")
	}
	userCohort, err := query.GetUserCohort(c.Db, c.SessionData.UserId)
	if err != nil {
		// TODO: Should probably check what the errors here are. Right now assume that cohort does not
		// exist
	}

	userModel := api.MyProfileResponse{
		UserPersonalInfo: api.UserPersonalInfo{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Gender:    user.Gender,
			Birthdate: user.Birthdate.Unix(),
		},
		UserContactInfo: api.UserContactInfo{
			Email: user.Email,
		},
	}

	if user.ExternalAuthData != nil {
		userModel.UserContactInfo.PhoneNumber = user.ExternalAuthData.PhoneNumber
		userModel.UserContactInfo.FbId = user.ExternalAuthData.FbUserId
	}

	if userCohort != nil {
		userModel.Cohort.CohortId = userCohort.CohortId
		userModel.Cohort.ProgramId = userCohort.ProgramId
		userModel.Cohort.GradYear = userCohort.GradYear
		userModel.Cohort.SequenceId = userCohort.SequenceId
	}

	c.Result = userModel
	return nil
}
