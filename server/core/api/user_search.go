package api

import (
	"letstalk/server/data"
)

type CommonUserSearchRequest struct {
	Size int `json:"size" binding:"required"`
}

type CohortUserSearchRequest struct {
	CommonUserSearchRequest
	CohortId data.TCohortID `json:"cohortId" binding:"required"`
}

type SimpleTraitUserSearchRequest struct {
	CommonUserSearchRequest
	SimpleTraitId data.TSimpleTraitID `json:"simpleTraitId" binding:"required"`
}

type PositionUserSearchRequest struct {
	CommonUserSearchRequest
	RoleId         data.TRoleID         `json:"roleId" binding:"required"`
	OrganizationId data.TOrganizationID `json:"organizationId" binding:"required"`
}

type GroupUserSearchRequest struct {
	CommonUserSearchRequest
	GroupId data.TGroupID `json:"groupId" binding:"required"`
}

type UserSearchResult struct {
	UserId     data.TUserID  `json:"userId"`
	FirstName  string        `json:"firstName"`
	LastName   string        `json:"lastName"`
	Gender     data.GenderID `json:"gender"`
	Cohort     *CohortV2     `json:"cohort"`
	ProfilePic *string       `json:"profilePic"`
	Reason     *string       `json:"reason"` // optional reason for the result being shown
}

// `isAnonymous` will be true if the searched term is sensitive and we don't want to be actually
// showing any results.
type UserSearchResponse struct {
	IsAnonymous bool               `json:"isAnonymous"`
	NumResults  int                `json:"numResults"` // Number of results, even if anonymous
	Results     []UserSearchResult `json:"results"`    // Empty if anonymous
}
