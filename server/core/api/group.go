package api

import "letstalk/server/data"

type UserGroup struct {
	Id        data.TUserGroupID `json:"id"`
	GroupId   data.TGroupID     `json:"groupId"`
	GroupName string            `json:"groupName"`
}
