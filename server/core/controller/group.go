package controller

import (
	"fmt"
	"letstalk/server/data"

	"letstalk/server/core/api"
	"letstalk/server/core/converters"
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
		return errs.NewRequestError(fmt.Sprintf("No survey for the %s group", req.GroupId))
	}
	survey, dbErr := query.GetSurvey(c.Db, c.SessionData.UserId, groupSurvey.Group)
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	userGroup, err := query.AddUserGroup(c.Db, c.SessionData.UserId, req.GroupId)
	if err != nil {
		return err
	}
	c.Result = &api.UserGroupSurvey{
		UserGroup: api.UserGroup{
			Id:      userGroup.Id,
			GroupId: userGroup.GroupId,
		},
		Survey: *survey,
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

// CreateManagedGroupController Create a new managed group, should require admin permission
func CreateManagedGroupController(c *ctx.Context) errs.Error {
	var req api.CreateGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	_, err := query.CreateManagedGroup(c.Db, c.SessionData.UserId, req.GroupName)
	if err != nil {
		return err
	}

	return nil
}

func AddAdminToManagedGroupController(c *ctx.Context) errs.Error {
	var req api.AddAdminToManagedGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	managesGroup, err := query.CheckAdminManagesGroup(c.Db, c.SessionData.UserId, req.GroupUUID)

	if err != nil {
		return err
	}

	if !managesGroup {
		return errs.NewForbiddenError("You are not allowed to make modifications to this group")
	}

	return query.AddAdminToManagedGroup(c.Db, req.AdminId, req.GroupUUID)
}

// GetAdminManagedGroupsController Get all groups that this admin manages
func GetAdminManagedGroupsController(c *ctx.Context) errs.Error {
	groups, err := query.GetManagedGroups(c.Db, c.SessionData.UserId)
	if err != nil {
		return err
	}
	var res api.GetAdminMangedGroupsResponse
	for _, group := range groups {
		res.ManagedGroups = append(res.ManagedGroups, api.AdminManagedGroup{
			GroupId:                   group.Group.GroupId,
			GroupName:                 group.Group.GroupName,
			ManagedGroupReferralEmail: converters.GetManagedGroupReferralLink(group.Group.GroupId),
		})
	}
	c.Result = res
	return nil
}

func RemoveUserManagedGroupController(c *ctx.Context) errs.Error {
	// make sure admin can touch this group
	var req api.RemoveUserGroupRequest2
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}
	adminManagesGroup, _ := query.CheckAdminManagesGroup(c.Db, c.SessionData.UserId, req.GroupId)
	if !adminManagesGroup {
		return errs.NewForbiddenError("You are not allowed to make modifications to that group")
	}
	return query.RemoveUserFromGroup(c.Db, req.UserId, req.GroupId)
}

// EnrollUserManagedGroupController Lets a user enroll themselves in a group
func EnrollUserManagedGroupController(c *ctx.Context) errs.Error {
	var req api.EnrollManagedGroupRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	return query.EnrollUserInManagedGroup(c.Db, c.SessionData.UserId, data.TGroupID(req.GroupUUID))
}

// EnrollUserInManagedGroupController Lets an admin enroll a user in another group
// TODO: do an access control check
func EnrollUserInManagedGroupController(c *ctx.Context) errs.Error {
	var req api.EnrollUserInManagedGroupRequest
	if err := c.GinContext.Bind(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	managesGroup, err := query.CheckAdminManagesGroup(c.Db, c.SessionData.UserId, req.GroupUUID)

	if err != nil {
		return err
	}

	if !managesGroup {
		return errs.NewForbiddenError("You are not allowed to make modifications to this group")
	}

	return query.EnrollUserInManagedGroup(c.Db, req.UserId, req.GroupUUID)
}

// EnrollUserInManagedGroupByEmailController Lets an admin enroll a user in another group
// TODO: do an access control check
func EnrollUserInManagedGroupByEmailController(c *ctx.Context) errs.Error {
	var req api.EnrollUserInManagedGroupByEmailRequest
	if err := c.GinContext.Bind(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	managesGroup, err := query.CheckAdminManagesGroup(c.Db, c.SessionData.UserId, req.GroupUUID)

	if err != nil {
		return err
	}

	if !managesGroup {
		return errs.NewForbiddenError("You are not allowed to make modifications to this group")
	}

	return query.EnrollUserInManagedGroupByEmail(c.Db, req.Email, req.GroupUUID)
}
