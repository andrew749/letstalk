package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

// Nuke user deletes all data about the given user. User caution with this endpoint as it has
// some effects.
func NukeUser(c *ctx.Context) errs.Error {
	var req api.NukeUserRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	if err := query.NukeUser(c.Db, req.Email, req.FirstName, req.LastName, req.UserId); err != nil {
		return errs.NewDbError(err)
	}
	return nil
}
