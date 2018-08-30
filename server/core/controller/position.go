package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func AddUserPositionController(c *ctx.Context) errs.Error {
	var req api.AddUserPositionRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	userPosition, err := query.AddUserPosition(
		c.Db,
		c.Es,
		c.SessionData.UserId,
		req.RoleId,
		req.RoleName,
		req.OrganizationId,
		req.OrganizationName,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return err
	}
	c.Result = &api.UserPosition{
		Id:               userPosition.Id,
		RoleId:           userPosition.RoleId,
		RoleName:         userPosition.RoleName,
		OrganizationId:   userPosition.OrganizationId,
		OrganizationName: userPosition.OrganizationName,
		OrganizationType: userPosition.OrganizationType,
		StartDate:        userPosition.StartDate,
		EndDate:          userPosition.EndDate,
	}
	return nil
}

func RemoveUserPositionController(c *ctx.Context) errs.Error {
	var req api.RemoveUserPositionRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.RemoveUserPosition(c.Db, c.SessionData.UserId, req.UserPositionId)
}
