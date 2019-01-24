package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/verify_link"
)

// Verifies link for a user
func VerifyLinkController(c *ctx.Context) errs.Error {
	var req api.VerifyLinkRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	if err := verify_link.ClickLink(c.Db, req.VerifyLinkId); err != nil {
		return err
	}

	c.Result = "Ok"
	return nil
}
