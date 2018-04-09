package request_to_match

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type GetCredentialOptionsResponse api.CredentialOptions

func GetCredentialOptionsController(c *ctx.Context) errs.Error {
	c.Result = api.GetCredentialOptions()
	return nil
}
