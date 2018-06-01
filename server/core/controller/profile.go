package controller

import (
	"strconv"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func ProfileEditController(c *ctx.Context) errs.Error {
	var request api.ProfileEditRequest
	err := c.GinContext.BindJSON(&request)
	if err != nil {
		return errs.NewClientError(err.Error())
	}

	if err := query.UpdateProfile(c.Db, c.SessionData.UserId, request); err != nil {
		return err
	}
	return nil
}

func GetMyProfileController(c *ctx.Context) errs.Error {
	userModel, err := query.GetProfile(c.Db, c.SessionData.UserId)
	if err != nil {
		return nil
	}

	c.Result = *userModel
	return nil
}

func GetMatchProfileController(c *ctx.Context) errs.Error {
	matchUserIdStr := c.GinContext.Param("userId")
	matchUserId, err := strconv.Atoi(matchUserIdStr)
	if err != nil {
		return errs.NewClientError(err.Error())
	}

	userModel, err := query.GetMatchProfile(c.Db, c.SessionData.UserId, matchUserId)
	if err != nil {
		return nil
	}

	c.Result = *userModel
	return nil
}
