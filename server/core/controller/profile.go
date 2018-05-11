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

	userModel := api.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
		Birthday:  user.Birthdate.Unix(),
		Email:     user.Email,
	}

	if user.ExternalAuthData != nil {
		userModel.PhoneNumber = user.ExternalAuthData.PhoneNumber
		userModel.FbId = user.ExternalAuthData.FbUserId
	}

	c.Result = userModel
	return nil
}
