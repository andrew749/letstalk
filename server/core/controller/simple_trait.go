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
	if err := query.AddUserSimpleTraitByName(
		c.Db,
		c.SessionData.UserId,
		req.SimpleTraitName,
	); err != nil {
		return err
	}
	return nil
}

func AddUserSimpleTraitByIdController(c *ctx.Context) errs.Error {
	var req api.AddUserSimpleTraitByIdRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	if err := query.AddUserSimpleTraitById(
		c.Db,
		c.SessionData.UserId,
		req.SimpleTraitId,
	); err != nil {
		return err
	}
	return nil
}

func RemoveUserSimpleTraitController(c *ctx.Context) errs.Error {
	var req api.RemoveUserSimpleTraitRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	if err := query.RemoveUserSimpleTrait(c.Db, c.SessionData.UserId, req.SimpleTraitId); err != nil {
		return err
	}
	return nil
}
