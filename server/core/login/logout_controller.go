package login

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
)

type LogoutResponse struct {
}

func LogoutHandler(c *ctx.Context) errs.Error {
	if c.SessionData.SessionId == nil {
		return errs.NewInternalError("Bad session token.")
	}

	// remove the session from list of active session
	if err := c.Db.Where("session_id = ?", c.SessionData.SessionId).Delete(data.Session{}).Error; err != nil {
		return errs.NewInternalError("Unable to delete session")
	}

	c.Result = "Ok"

	return nil
}
