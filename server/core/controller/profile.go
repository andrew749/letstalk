package controller

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

func ProfileEditController(c *ctx.Context) errs.Error {
	var request api.ProfileEditRequest
	err := c.GinContext.BindJSON(&request)
	if err != nil {
		return errs.NewClientError(err.Error())
	}

	bday := time.Unix(request.Birthdate, 0)

	tx := c.Db.Begin()
	err = tx.Model(&data.User{}).Where(&data.User{
		UserId: c.SessionData.UserId,
	}).Update(data.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Gender:    request.Gender,
		Birthdate: &bday,
	}).Error
	if err != nil {
		tx.Rollback()
		return errs.NewDbError(err)
	}

	err = tx.Model(&data.UserCohort{}).Where(&data.UserCohort{
		UserId: c.SessionData.UserId,
	}).Update(data.UserCohort{CohortId: request.CohortId}).Error
	if err != nil {
		tx.Rollback()
		return errs.NewDbError(err)
	}

	if request.PhoneNumber != nil {
		err = tx.Model(&data.ExternalAuthData{}).Where(&data.ExternalAuthData{
			UserId: c.SessionData.UserId,
		}).Update(data.ExternalAuthData{
			PhoneNumber: request.PhoneNumber,
		}).FirstOrCreate(&data.ExternalAuthData{}).Error
		if err != nil {
			tx.Rollback()
			return errs.NewDbError(err)
		}
	}

	tx.Commit()
	return nil
}

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
