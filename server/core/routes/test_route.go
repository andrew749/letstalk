package routes

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

func GetTest(c *ctx.Context) errs.Error {
	result := struct{ Response string }{"test controller"}
	c.Result = result
	return nil
}

func GetTestAuth(c *ctx.Context) errs.Error {
	c.Result = struct{ Message string }{"test"}
	return nil
}
