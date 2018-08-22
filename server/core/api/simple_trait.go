package api

import "letstalk/server/data"

type SimpleTrait struct {
	Id          data.TSimpleTraitID  `json:"id"`
	Name        string               `json:"name"`
	Type        data.SimpleTraitType `json:"type"`
	IsSensitive bool                 `json:"isSensitive"`
}

type AddUserSimpleTraitByNameRequest struct {
	SimpleTraitName string `json:"name"`
}

type AddUserSimpleTraitByIdRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId"`
}

type RemoveUserSimpleTraitRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId"`
}
