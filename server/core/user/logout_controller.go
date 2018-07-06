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

	// TODO: Refactor this to not live in the controller
	// remove the session from list of active session
	err := c.Db.Where("session_id = ?", c.SessionData.SessionId).Delete(data.Session{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}

	err = c.Db.Where("session_id = ?", c.SessionData.SessionId).Delete(data.NotificationToken{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}

	c.Result = "Ok"

	return nil
}
