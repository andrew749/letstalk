package bootstrap

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

func convertUserToRelationshipDataModel(
	user data.User,
	description *string,
	userType api.UserType,
) *api.BootstrapUserRelationshipDataModel {
	var (
		fbId        *string
		phoneNumber *string
		cohort      *api.Cohort
	)
	if user.ExternalAuthData != nil {
		fbId = user.ExternalAuthData.FbUserId
		phoneNumber = user.ExternalAuthData.PhoneNumber
	}

	if user.Cohort != nil && user.Cohort.Cohort != nil {
		cohort = &api.Cohort{
			CohortId:   user.Cohort.Cohort.CohortId,
			ProgramId:  user.Cohort.Cohort.ProgramId,
			SequenceId: user.Cohort.Cohort.SequenceId,
			GradYear:   user.Cohort.Cohort.GradYear,
		}
	}

	return &api.BootstrapUserRelationshipDataModel{
		UserId:      user.UserId,
		UserType:    userType,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		FbId:        fbId,
		PhoneNumber: phoneNumber,
		Description: description,
		Cohort:      cohort,
	}
}

func getDescriptionForRequestToMatch(requestMatching data.RequestMatching) *string {
	if requestMatching.Credential != nil {
		return &requestMatching.Credential.Name
	}
	return nil
}

/**
 * Returns what the current status of a user is
 */
func GetCurrentUserBoostrapStatusController(c *ctx.Context) errs.Error {
	var (
		err      error
		response = api.BootstrapResponse{State: api.ACCOUNT_CREATED}
		userId   = c.SessionData.UserId
	)

	onboardingInfo, err := onboarding.GetOnboardingInfo(c.Db, userId)
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
	mentors, err := query.GetMentorsByMenteeId(c.Db, userId) // Matchings where user is the mentee.
	if err != nil {
		return errs.NewDbError(err)
	}
	mentees, err := query.GetMenteesByMentorId(c.Db, userId) // Matchings where user is the mentor.
	if err != nil {
		return errs.NewDbError(err)
	}
	askers, err := query.GetAskersByAnswererId(c.Db, userId) // Request matchings where user is answerer.
	if err != nil {
		return errs.NewDbError(err)
	}
	answerers, err := query.GetAnswerersByAskerId(c.Db, userId) // Request matchings where user is asker.
	if err != nil {
		return errs.NewDbError(err)
	}

	// Construct relationship api objects.
	relationships := make(
		[]*api.BootstrapUserRelationshipDataModel,
		0,
		len(mentors)+len(mentees)+len(askers)+len(answerers),
	)
	for _, mentor := range mentors {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(*mentor.MentorUser, nil, api.USER_TYPE_MENTOR),
		)
	}
	for _, mentee := range mentees {
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(*mentee.MenteeUser, nil, api.USER_TYPE_MENTEE),
		)
	}
	for _, asker := range askers {
		description := getDescriptionForRequestToMatch(asker)
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(
				*asker.AskerUser,
				description,
				api.USER_TYPE_ASKER,
			),
		)
	}
	for _, answerer := range answerers {
		description := getDescriptionForRequestToMatch(answerer)
		relationships = append(
			relationships,
			convertUserToRelationshipDataModel(
				*answerer.AnswererUser,
				description,
				api.USER_TYPE_ANSWERER,
			),
		)
	}
	if len(relationships) > 0 {
		response.State = api.ACCOUNT_MATCHED
	}

	response.Relationships = relationships
	c.Result = response

	return nil
}
