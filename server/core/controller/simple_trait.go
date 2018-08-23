package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func AddUserSimpleTraitByNameController(c *ctx.Context) errs.Error {
	var req api.AddUserSimpleTraitByNameRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.AddUserSimpleTraitByName(
		c.Db,
		c.Es,
		c.SessionData.UserId,
		req.SimpleTraitName,
	)
}

func AddUserSimpleTraitByIdController(c *ctx.Context) errs.Error {
	var req api.AddUserSimpleTraitByIdRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.AddUserSimpleTraitById(
		c.Db,
		c.SessionData.UserId,
		req.SimpleTraitId,
	)
}

func RemoveUserSimpleTraitController(c *ctx.Context) errs.Error {
	var req api.RemoveUserSimpleTraitRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.RemoveUserSimpleTrait(c.Db, c.SessionData.UserId, req.SimpleTraitId)
}
