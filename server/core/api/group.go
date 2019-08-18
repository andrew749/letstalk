package api

import "letstalk/server/data"

type UserGroup struct {
	Id        data.TUserGroupID `json:"id"`
	GroupId   data.TGroupID     `json:"groupId"`
	GroupName string            `json:"groupName"`
}

type AddUserGroupRequest struct {
	GroupId data.TGroupID `json:"groupId"`
}

type RemoveUserGroupRequest struct {
	UserGroupId data.TUserGroupID `json:"userGroupId"`
}

// Another way to remove a user from the group
type RemoveUserGroupRequest2 struct {
	UserId  data.TUserID  `json:"userId"`
	GroupId data.TGroupID `json:"groupId"`
}

type CreateGroupRequest struct {
	GroupName string `json:"groupName"`
}

type GetAdminMangedGroupsResponse struct {
	ManagedGroups []AdminManagedGroup `json:"managedGroups"`
}

type AdminManagedGroup struct {
	GroupId                   data.TGroupID `json:"groupId"`
	GroupName                 string        `json:"groupName"`
	ManagedGroupReferralEmail string        `json:"managedGroupReferralEmail"`
}

type EnrollManagedGroupRequest struct {
	GroupUUID string `json:"groupUUID"`
}

type EnrollUserInManagedGroupRequest struct {
	GroupUUID data.TGroupID `json:"groupUUID"`
	UserId    data.TUserID  `json:"userId"`
}

type EnrollUserInManagedGroupByEmailRequest struct {
	GroupUUID data.TGroupID `json:"groupUUID"`
	Email     string        `json:"email"`
}

type GroupMemberStatus string

const (
	GROUP_MEMBER_STATUS_SIGNED_UP GroupMemberStatus = "SIGNED_UP"
	GROUP_MEMBER_STATUS_ONBOARDED GroupMemberStatus = "ONBOARDED"
	GROUP_MEMBER_STATUS_MATCHED   GroupMemberStatus = "MATCHED"
)

type GroupMember struct {
	User   UserPersonalInfo  `json:"user" binding:"required"`
	Email  string            `json:"email" binding:"required"`
	Status GroupMemberStatus `json:"status" binding:"required"`
	Cohort *CohortV2         `json:"cohort"`
}
