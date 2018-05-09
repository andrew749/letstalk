package api

import (
	"letstalk/server/data"
)

type BootstrapState string

/**
 * These states will likely change.
 * Current a later state implies that the previous states are satisfied
 * This is currently a linear state hierarchy
 */
const (
	ACCOUNT_CREATED BootstrapState = "account_created" // first state
	ACCOUNT_SETUP   BootstrapState = "account_setup"   // the account has enough information to proceed
	ACCOUNT_MATCHED BootstrapState = "account_matched" // account has been matched a peer
)

type BootstrapUserRelationshipDataModel struct {
	User      int      `json:"userId" binding:"required"`
	UserType  UserType `json:"userType" binding:"required"`
	FirstName string   `json:"firstName" binding:"required"`
	LastName  string   `json:"lastName" binding:"required"`
	Email     string   `json:"email" binding:"required"`
	FbId      *string  `json:"fbId"`
}

type BootstrapResponse struct {
	State            BootstrapState                        `json:"state" binding:"required"`
	Relationships    []*BootstrapUserRelationshipDataModel `json:"relationships" binding:"required"`
	Cohort           *data.Cohort                          `json:"cohort" binding:"required"`
	Me               *data.User                            `json:"me" binding:"required"`
	OnboardingStatus *OnboardingStatus                     `json:"onboardingStatus" binding:"required"`
}
