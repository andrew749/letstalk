package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

func AddUserGroupController(c *ctx.Context) errs.Error {
	var req api.AddUserGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	userGroup, err := query.AddUserGroup(c.Db, c.SessionData.UserId, req.GroupId, req.GroupName)
	if err != nil {
		return err
	}
	c.Result = &api.UserGroup{
		Id:        userGroup.Id,
		GroupId:   userGroup.GroupId,
		GroupName: userGroup.GroupName,
	}
	return nil
}

func RemoveUserGroupController(c *ctx.Context) errs.Error {
	var req api.RemoveUserGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	return query.RemoveUserGroup(c.Db, c.SessionData.UserId, req.UserGroupId)
}
