package api

import "letstalk/server/data"

type SimpleTrait struct {
	Id          data.TSimpleTraitID  `json:"id" binding:"required"`
	Name        string               `json:"name" binding:"required"`
	Type        data.SimpleTraitType `json:"type" binding:"required"`
	IsSensitive bool                 `json:"isSensitive" binding:"required"`
}

type UserSimpleTrait struct {
	Id                     data.TUserSimpleTraitID `json:"id"`
	SimpleTraitId          data.TSimpleTraitID     `json:"simpleTraitId"`
	SimpleTraitName        string                  `json:"simpleTraitName"`
	SimpleTraitType        data.SimpleTraitType    `json:"simpleTraitType"`
	SimpleTraitIsSensitive bool                    `json:"simpleTraitIsSensitive"`
}

type AddUserSimpleTraitByNameRequest struct {
	SimpleTraitName string `json:"name" binding:"required"`
}

type AddUserSimpleTraitByIdRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId" binding:"required"`
}

type RemoveUserSimpleTraitRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId" binding:"required"`
}
