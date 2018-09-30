package user

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
)

func LogoutHandler(c *ctx.Context) errs.Error {
	if c.SessionData.SessionId == nil {
		return errs.NewInternalError("Bad session token.")
	}

	if err := data.DeleteSession(c.Db, *c.SessionData.SessionId); err != nil {
		return errs.NewDbError(err)
	}

	c.Result = "Ok"

	return nil
}
