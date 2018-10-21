package bootstrap

import (
	"github.com/romana/rlog"
	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/core/survey"
)

func convertUserToRelationshipDataModel(
	user data.User,
	matchingState data.MatchingState,
	description *string,
	userType api.UserType,
) *api.BootstrapUserRelationshipDataModel {
	var (
		fbId        *string
		fbLink      *string
		phoneNumber *string
		cohort      *api.Cohort
	)
	if user.ExternalAuthData != nil {
		fbId = user.ExternalAuthData.FbUserId
		fbLink = user.ExternalAuthData.FbProfileLink
		phoneNumber = user.ExternalAuthData.PhoneNumber
	}

	if user.Cohort != nil && user.Cohort.Cohort != nil {
		// NOTE: New API will allow for null sequence ids.
		sequenceId := ""
		if user.Cohort.Cohort.SequenceId != nil {
			sequenceId = *user.Cohort.Cohort.SequenceId
		}

		cohort = &api.Cohort{
			CohortId:   user.Cohort.Cohort.CohortId,
			ProgramId:  user.Cohort.Cohort.ProgramId,
			SequenceId: sequenceId,
			GradYear:   user.Cohort.Cohort.GradYear,
		}
	}

	return &api.BootstrapUserRelationshipDataModel{
		UserId:        user.UserId,
		UserType:      userType,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		FbId:          fbId,
		FBLink:        fbLink,
		PhoneNumber:   phoneNumber,
		Description:   description,
		Cohort:        cohort,
		MatchingState: matchingState,
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

	user, err := query.GetUserById(c.Db, userId)
	if err != nil {
		return errs.NewInternalError("Authenticated user not found")
	}
	if !user.IsEmailVerified {
		// User email not yet verified, don't proceed to onboarding.
		c.Result = response
		return nil
	} else {
		response.State = api.ACCOUNT_EMAIL_VERIFIED
	}

	onboardingInfo, err := onboarding.GetOnboardingInfo(c.Db, userId)
	if err != nil {
		return errs.NewDbError(err)
	}
	response.Cohort = onboardingInfo.UserCohort
	response.OnboardingStatus = &api.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}

	// Fetch user's survey information
	surveyResponses, surveyErr := survey.GetSurveyResponses(c.Db, user.UserId, survey.Generic_v1.Group)
	if surveyErr != nil {
		return surveyErr
	}
	userSurvey := survey.Generic_v1
	if surveyResponses != nil {
		userSurvey.Responses = surveyResponses
	}
	response.Survey = &userSurvey

	if onboardingInfo.State != api.ONBOARDING_DONE {
		// Onboarding not done. We don't need to get connections.
		c.Result = response
		return nil
	} else {
		response.State = api.ACCOUNT_SETUP
	}

	// Fetch all user's connections.

	connections, err := query.GetAllConnections(c.Db, userId)
	if err != nil {
		return errs.NewDbError(err)
	}

	if len(connections) > 0 {
		response.State = api.ACCOUNT_MATCHED
	}
	response.Connections.IncomingRequests = make([]*api.ConnectionRequestWithName, 0)
	response.Connections.OutgoingRequests = make([]*api.ConnectionRequestWithName, 0)
	response.Connections.Mentees = make([]*api.BootstrapConnection, 0)
	response.Connections.Mentors = make([]*api.BootstrapConnection, 0)
	response.Connections.Peers = make([]*api.BootstrapConnection, 0)

	for _, conn := range connections {
		connUserId := conn.UserOneId
		connUser := conn.UserOne
		if conn.UserOneId == user.UserId {
			connUserId = conn.UserTwoId
			connUser = conn.UserTwo
		}
		if connUser == nil {
			return errs.NewInternalError("Failed to load connection user data")
		}
		if conn.AcceptedAt == nil {
			if conn.UserOneId == user.UserId {
				// Auth user is the requestor.
				connApi := api.ConnectionRequestWithName{
					ConnectionRequest: connection.DataToApi(connUserId, conn),
					FirstName: conn.UserTwo.FirstName,
					LastName: conn.UserTwo.LastName,
				}
				rlog.Debug("adding outgoing request", connApi)
				response.Connections.OutgoingRequests =
					append(response.Connections.OutgoingRequests, &connApi)
			} else {
				// Auth user is the requestee.
				connApi := api.ConnectionRequestWithName{
					ConnectionRequest: connection.DataToApi(connUserId, conn),
					FirstName: conn.UserOne.FirstName,
					LastName: conn.UserOne.LastName,
				}
				rlog.Debug("adding incoming request", connApi)
				response.Connections.IncomingRequests =
					append(response.Connections.IncomingRequests, &connApi)
			}
		} else {
			if conn.Mentorship != nil {
				// Connection has been upgraded to a mentorship.
				if conn.Mentorship.MentorUserId == user.UserId {
					// Auth user is the mentor.
					bc := api.BootstrapConnection{
						Request: connection.DataToApi(connUserId, conn),
						UserProfile: *convertUserToRelationshipDataModel(
							*connUser,
							data.MATCHING_STATE_UNKNOWN,
							conn.Intent.SearchedTrait,
							api.USER_TYPE_MENTEE,
						),
					}
					rlog.Debug("adding mentee", bc)
					response.Connections.Mentees =
						append(response.Connections.Mentees, &bc)
				} else {
					// Auth user is the mentee.
					bc := api.BootstrapConnection{
						Request: connection.DataToApi(connUserId, conn),
						UserProfile: *convertUserToRelationshipDataModel(
							*connUser,
							data.MATCHING_STATE_UNKNOWN,
							conn.Intent.SearchedTrait,
							api.USER_TYPE_MENTOR,
						),
					}
					rlog.Debug("adding mentor", bc)
					response.Connections.Mentors =
						append(response.Connections.Mentors, &bc)
				}
			} else {
				userType := api.USER_TYPE_ASKER
				if conn.UserOneId == user.UserId {
					userType = api.USER_TYPE_ANSWERER
				}
				// Connection is not a mentorship.
				bc := api.BootstrapConnection{
					Request: connection.DataToApi(connUserId, conn),
					UserProfile: *convertUserToRelationshipDataModel(
						*connUser,
						data.MATCHING_STATE_UNKNOWN,
						conn.Intent.SearchedTrait,
						userType,
					),
				}
				rlog.Debug("adding other connection", bc)
				response.Connections.Peers =
					append(response.Connections.Peers, &bc)
			}
		}
	}

	c.Result = response
	return nil
}
