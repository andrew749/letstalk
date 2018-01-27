package routes

import (
	"letstalk/server/core/ctx"
)

func GetTest(c *ctx.Context) {
	result := struct{ Response string }{ "test controller" }
	c.Result = result
}
