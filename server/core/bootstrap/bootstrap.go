package bootstrap

import (
	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/user_state"
	"letstalk/server/data"

	"github.com/romana/rlog"
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
		err      errs.Error
		response = api.BootstrapResponse{State: api.ACCOUNT_CREATED}
		userId   = c.SessionData.UserId
	)

	state, err := user_state.GetUserState(c.Db, userId)
	if err != nil {
		return err
	}
	response.State = *state
	if *state != api.ACCOUNT_SETUP {
		// Account isn't setup yet, so we shouldn't be fetching connections, even if they have them
		c.Result = response
		return nil
	}

	// Fetch all user's connections.

	connections, dbErr := query.GetAllConnections(c.Db, userId)
	if dbErr != nil {
		return errs.NewDbError(dbErr)
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
		if conn.UserOneId == userId {
			connUserId = conn.UserTwoId
			connUser = conn.UserTwo
		}
		if connUser == nil {
			return errs.NewInternalError("Failed to load connection user data")
		}
		if conn.AcceptedAt == nil {
			connApi := api.ConnectionRequestWithName{
				ConnectionRequest: connection.DataToApi(connUserId, conn),
				FirstName:         connUser.FirstName,
				LastName:          connUser.LastName,
			}
			if conn.UserOneId == userId {
				// Auth user is the requestor.
				rlog.Debug("adding outgoing request", connApi)
				response.Connections.OutgoingRequests =
					append(response.Connections.OutgoingRequests, &connApi)
			} else {
				// Auth user is the requestee.
				rlog.Debug("adding incoming request", connApi)
				response.Connections.IncomingRequests =
					append(response.Connections.IncomingRequests, &connApi)
			}
		} else {
			if conn.Mentorship != nil {
				userType := api.USER_TYPE_MENTOR
				if conn.Mentorship.MentorUserId == userId {
					userType = api.USER_TYPE_MENTEE
				}

				bc := api.BootstrapConnection{
					Request: connection.DataToApi(connUserId, conn),
					UserProfile: *convertUserToRelationshipDataModel(
						*connUser,
						data.MATCHING_STATE_UNKNOWN,
						conn.Intent.SearchedTrait,
						userType,
					),
				}

				if conn.Mentorship.MentorUserId == userId {
					// Auth user is the mentor.
					rlog.Debug("adding mentee", bc)
					response.Connections.Mentees =
						append(response.Connections.Mentees, &bc)
				} else {
					// Auth user is the mentee.
					rlog.Debug("adding mentor", bc)
					response.Connections.Mentors =
						append(response.Connections.Mentors, &bc)
				}
			} else {
				userType := api.USER_TYPE_ASKER
				if conn.UserOneId == userId {
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
