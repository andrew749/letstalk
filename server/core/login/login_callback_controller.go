package login

import (
	"net/http"
	"letstalk/server/core/ctx"
)

type CallbackResponse struct {
	Status string
	Code   string
	State  string
}

func PostLoginSucceed(c *ctx.Context) {
	status := c.GinContext.Query("status")
	code := c.GinContext.Query("code")
	state := c.GinContext.Query("state")

	// authenticated data from the provider
	_ = CallbackResponse{status, code, state}
}
