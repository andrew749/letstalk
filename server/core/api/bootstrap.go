package api

import (
	"letstalk/server/data"
)

type BootstrapUserRelationshipDataModel struct {
	UserId        data.TUserID       `json:"userId" binding:"required"`
	UserType      UserType           `json:"userType" binding:"required"`
	FirstName     string             `json:"firstName" binding:"required"`
	LastName      string             `json:"lastName" binding:"required"`
	Email         string             `json:"email" binding:"required"`
	FbId          *string            `json:"fbId"`
	FBLink        *string            `json:"fbLink" binding:"required"`
	PhoneNumber   *string            `json:"phoneNumber"`
	Cohort        *Cohort            `json:"cohort"`
	Description   *string            `json:"description"`
	MatchingState data.MatchingState `json:"matchingState"`
}

type ConnectionRequestWithName struct {
	ConnectionRequest
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type BootstrapConnection struct {
	UserProfile BootstrapUserRelationshipDataModel `json:"userProfile" binding:"required"`
	Request     ConnectionRequest                  `json:"request" binding:"required"`
}

type BootstrapConnections struct {
	OutgoingRequests []*ConnectionRequestWithName `json:"outgoingRequests" binding:"required"`
	IncomingRequests []*ConnectionRequestWithName `json:"incomingRequests" binding:"required"`
	Mentors          []*BootstrapConnection       `json:"mentors" binding:"required"`
	Mentees          []*BootstrapConnection       `json:"mentees" binding:"required"`
	Peers            []*BootstrapConnection       `json:"peers" binding:"required"`
}

type BootstrapResponse struct {
	State       UserState            `json:"state" binding:"required"`
	Connections BootstrapConnections `json:"connections" binding:"required"`
}
