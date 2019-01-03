package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
)

func AddExpoDeviceToken(c *ctx.Context) errs.Error {
	var req api.AddExpoDeviceTokenRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	return handleAddExpoDeviceToken(c, req)
}

func handleAddExpoDeviceToken(c *ctx.Context, req api.AddExpoDeviceTokenRequest) errs.Error {
	err := data.AddExpoDeviceTokenForUser(c.Db, c.SessionData.UserId, req.Token)
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
