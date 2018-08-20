package api

import "letstalk/server/data"

type AddUserSimpleTraitByNameRequest struct {
	SimpleTraitName string `json:"name"`
}

type AddUserSimpleTraitByIdRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId"`
}

type RemoveUserSimpleTraitRequest struct {
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId"`
}
