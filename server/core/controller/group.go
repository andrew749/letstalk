package controller

import (
	"fmt"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/survey"
)

func AddUserGroupController(c *ctx.Context) errs.Error {
	var req api.AddUserGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	groupSurvey := survey.GetSurveyDefinitionByGroupId(req.GroupId)
	if groupSurvey == nil {
		return errs.NewRequestError(fmt.Sprintf("No survey for the %s group", req.GroupName))
	}
	userSurvey, dbErr := query.GetUserSurvey(c.Db, c.SessionData.UserId, groupSurvey.Group)
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	if userSurvey != nil {
		groupSurvey.Responses = &userSurvey.Responses
	}
	userGroup, err := query.AddUserGroup(c.Db, c.SessionData.UserId, req.GroupId, req.GroupName)
	if err != nil {
		return err
	}
	c.Result = &api.UserGroupSurvey{
		UserGroup: api.UserGroup{
			Id:        userGroup.Id,
			GroupId:   userGroup.GroupId,
			GroupName: userGroup.GroupName,
		},
		Survey: *groupSurvey,
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
