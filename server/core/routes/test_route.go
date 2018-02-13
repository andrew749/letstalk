package routes

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

func GetTest(c *ctx.Context) errs.Error {
	result := struct{ Response string `json:"response"` }{"test controller"}
	c.Result = result
	return nil
}

func GetTestAuth(c *ctx.Context) errs.Error {
	c.Result = struct{ Message string `json:"response"` }{"test"}
	return nil
}
