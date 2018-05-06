package bootstrap

import (
	"letstalk/server/core/query"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
	"letstalk/server/data"
	"letstalk/server/core/api"
)

func convertUserToRelationshipDataModel(user *data.User, isMentor bool) *api.BootstrapUserRelationshipDataModel {
	var userType api.UserType
	if isMentor == true {
		userType = api.USER_TYPE_MENTOR
	} else {
		userType = api.USER_TYPE_MENTEE
	}
	return &api.BootstrapUserRelationshipDataModel{
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
	user, err := query.GetUserById(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewInternalError("Unable to get user data.")
	}
	// since this method is authenticated the account needs to exist.
	var response = api.BootstrapResponse{
		State: api.ACCOUNT_CREATED,
		Me:    user,
	}

	onboardingInfo, err := onboarding.GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	response.Cohort = onboardingInfo.UserCohort
	response.OnboardingStatus = &api.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}

	if onboardingInfo.State != api.ONBOARDING_DONE {
		// Onboarding not done. We don't need to get relationships.
		c.Result = response
		return nil
	} else {
		response.State = api.ACCOUNT_SETUP
	}

	// Fetch mentors and mentees.
	mentors, err := query.GetMentorsByMenteeId(c.Db, user.UserId) // Matchings where user is the mentee.
	if err != nil {
		return errs.NewDbError(err)
	}
	mentees, err := query.GetMenteesByMentorId(c.Db, user.UserId) // Matchings where user is the mentor.
	if err != nil {
		return errs.NewDbError(err)
	}

	// Construct relationship api objects.
	relationships := make([]*api.BootstrapUserRelationshipDataModel, 0, len(mentors) + len(mentees))
	for _, mentor := range mentors {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(&mentor.MentorUser, true),
		)
	}
	for _, mentee := range mentees {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(&mentee.MenteeUser, false),
		)
	}
	if len(relationships) > 0 {
		response.State = api.ACCOUNT_MATCHED
	}

	response.Relationships = relationships
	c.Result = response

	return nil
}
