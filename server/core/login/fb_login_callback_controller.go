package login

/**
 * TODO(acod)
 * Massive WIP: do not touch
 */

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type CallbackResponse struct {
	Status string
	Code   string
	State  string
}

func GetLoginResponse(c *ctx.Context) errs.Error {
	status := c.GinContext.Query("status")
	code := c.GinContext.Query("code")
	state := c.GinContext.Query("state")

	// authenticated data from the provider
	_ = CallbackResponse{status, code, state}
	return nil
}
