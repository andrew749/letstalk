package login

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
)

type LogoutRequest struct {
}

type LogoutResponse struct {
}

func LogoutHandler(c *ctx.Context) errs.Error {
	var req LogoutRequest
	err := c.GinContext.BindJSON(&req)
	if err != nil {
		return errs.NewClientError("Bad logout request")
	}

	if c.SessionData.SessionId == nil {
		return errs.NewInternalError("Bad session token.")
	}

	// remove the session from list of active session
	_, err = data.SessionsObjs.
		Delete().
		Where(data.SessionsObjs.FilterSessionId("=", *c.SessionData.SessionId)).
		One(c.Db)

	if err != nil {
		return errs.NewInternalError("Unable to remove session")
	}

	return nil
}
