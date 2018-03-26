package bootstrap

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
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
	User      int          `json:"userId" binding:"required"`
	UserType  api.UserType `json:"userType" binding:"required"`
	FirstName string       `json:"firstName" binding:"required"`
	LastName  string       `json:"lastName" binding:"required"`
	Email     string       `json:"email" binding:"required"`
}

type BootstrapResponse struct {
	State            BootstrapState                        `json:"state" binding:"required"`
	Relationships    []*BootstrapUserRelationshipDataModel `json:"relationships" binding:"required"`
	Cohort           *data.Cohort                          `json:"cohort" binding:"required"`
	Me               *data.User                            `json:"me" binding:"required"`
	OnboardingStatus *onboarding.OnboardingStatus          `json:"onboardingStatus" binding:"required"`
}

func convertUserToRelationshipDataModel(user *data.User, isMentor bool) *BootstrapUserRelationshipDataModel {
	var userType api.UserType
	if isMentor == true {
		userType = api.USER_TYPE_MENTOR
	} else {
		userType = api.USER_TYPE_MENTEE
	}
	return &BootstrapUserRelationshipDataModel{
		User:      user.UserId,
		UserType:  userType,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

/**
 * Returns what the current status of a user is
 */
func GetCurrentUserBoostrapStatusController(c *ctx.Context) errs.Error {
	user, err := api.GetFullUserWithId(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewInternalError("Unable to get user data.")
	}
	// since this method is authenticated the account needs to exist.
	var response = BootstrapResponse{
		State: ACCOUNT_CREATED,
		Me:    user,
	}

	onboardingInfo, err := onboarding.GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	response.Cohort = onboardingInfo.UserCohort
	response.OnboardingStatus = &onboarding.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}

	if onboardingInfo.State != onboarding.ONBOARDING_DONE {
		// Onboarding not done. We don't need to get relationships.
		c.Result = response
		return nil
	} else {
		response.State = ACCOUNT_SETUP
	}

	if len(user.Mentors) > 0 {
		response.State = ACCOUNT_MATCHED
	}

	relationships := make([]*BootstrapUserRelationshipDataModel, 0)
	// get all mentors
	for _, mentor := range user.Mentors {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(mentor, true),
		)
	}

	// get all mentees
	for _, mentee := range user.Mentees {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(mentee, false),
		)
	}

	response.Relationships = relationships
	c.Result = response

	return nil
}
