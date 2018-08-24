package api

import "letstalk/server/data"

type Role struct {
	Id   data.TRoleID `json:"id" binding:"required"`
	Name string       `json:"name" binding:"required"`
}

type AddUserPositionRequest struct {
	RoleId           *data.TRoleID         `json:"roleId"`
	RoleName         *string               `json:"roleName"`
	OrganizationId   *data.TOrganizationID `json:"organizationId"`
	OrganizationName *string               `json:"organizationName"`
	StartDate        string                `json:"startDate" binding:"required"`
	EndDate          *string               `json:"startDate"`
}

type RemoveUserPositionRequest struct {
	UserPositionId data.TUserPositionID `json:"userPositionId" binding:"required"`
}
