package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

func dataToApiSimpleTrait(userTrait *data.UserSimpleTrait) *api.UserSimpleTrait {
	return &api.UserSimpleTrait{
		Id:                     userTrait.Id,
		SimpleTraitId:          userTrait.SimpleTraitId,
		SimpleTraitName:        userTrait.SimpleTraitName,
		SimpleTraitType:        userTrait.SimpleTraitType,
		SimpleTraitIsSensitive: userTrait.SimpleTraitIsSensitive,
	}
}

func AddUserSimpleTraitByNameController(c *ctx.Context) errs.Error {
	var req api.AddUserSimpleTraitByNameRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	userTrait, err := query.AddUserSimpleTraitByName(
		c.Db,
		c.Es,
		c.SessionData.UserId,
		req.SimpleTraitName,
	)
	if err != nil {
		return err
	}
	c.Result = dataToApiSimpleTrait(userTrait)
	return nil
}

func AddUserSimpleTraitByIdController(c *ctx.Context) errs.Error {
	var req api.AddUserSimpleTraitByIdRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	userTrait, err := query.AddUserSimpleTraitById(
		c.Db,
		c.SessionData.UserId,
		req.SimpleTraitId,
	)
	if err != nil {
		return err
	}
	c.Result = dataToApiSimpleTrait(userTrait)
	return nil
}

func RemoveUserSimpleTraitController(c *ctx.Context) errs.Error {
	var req api.RemoveUserSimpleTraitRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.RemoveUserSimpleTrait(c.Db, c.SessionData.UserId, req.UserSimpleTraitId)
}
