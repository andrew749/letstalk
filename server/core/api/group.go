package api

import "letstalk/server/data"

type UserGroup struct {
	Id        data.TUserGroupID `json:"id"`
	GroupId   data.TGroupID     `json:"groupId"`
	GroupName string            `json:"groupName"`
}

type AddUserGroupRequest struct {
	GroupId   data.TGroupID `json:"groupId"`
	GroupName string        `json:"groupName"`
}

type RemoveUserGroupRequest struct {
	UserGroupId data.TUserGroupID `json:"userGroupId"`
}

type CreateGroupRequest struct {
	GroupName string `json:"groupName"`
}

type GetAdminMangedGroupsResponse struct {
	ManagedGroups []AdminManagedGroup `json:"managedGroups"`
}

type AdminManagedGroup struct {
	GroupName                 string `json:"groupName"`
	ManagedGroupReferralEmail string `json:"managedGroupReferralEmail"`
}

type EnrollManagedGroupRequest struct {
	GroupUUID string `json:"groupUUID"`
}
