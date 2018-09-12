package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func SimpleTraitUserSearchController(c *ctx.Context) errs.Error {
	var req api.SimpleTraitUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersBySimpleTrait(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}

func PositionUserSearchController(c *ctx.Context) errs.Error {
	var req api.PositionUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersByPosition(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}

func CohortUserSearchController(c *ctx.Context) errs.Error {
	var req api.CohortUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersByCohort(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}
