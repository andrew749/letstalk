package bootstrap

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
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

type UserType string

// the roles a user can take in a relationship
const (
	MENTOR UserType = "mentor"
	MENTEE UserType = "mentee"
)

type BootstrapUserRelationshipDataModel struct {
	User     int      `json:"user_id" binding:"required"`
	UserType UserType `json:"user_type" binding:"required"`
}

type BootstrapResponse struct {
	State           BootstrapState                       `json:"state" binding:"required"`
	Relatationships []BootstrapUserRelationshipDataModel `json:"relationships" binding:"required"`
	Cohort          *data.Cohort                         `json:"cohort" binding:"required"`
}

func convertMatchingToRelationshipDataModel(
	userId int,
	userType UserType,
) BootstrapUserRelationshipDataModel {
	return BootstrapUserRelationshipDataModel{
		User:     userId,
		UserType: userType,
	}
}

/**
 * Returns what the current status of a user is
 */
func GetCurrentUserBoostrapStatusController(c *ctx.Context) errs.Error {
	// since this method is authenticated the account needs to exist.
	var response = BootstrapResponse{
		State: ACCOUNT_CREATED,
	}

	// check if the account has been onboarded
	userCohort, err := api.GetUserCohort(c.Db, c.SessionData.UserId)

	if err == nil {
		response.State = ACCOUNT_SETUP
		response.Cohort = userCohort
	}

	relationships := make([]BootstrapUserRelationshipDataModel, 0)

	// find any of this person's mentees
	mentee_matchings, err := data.MatchingsObjs.
		Select().
		Where(data.MatchingsObjs.FilterMentor("=", c.SessionData.UserId)).
		List(c.Db)

	if err == nil && len(mentee_matchings) > 0 {
		response.State = ACCOUNT_MATCHED
		// create array of matchings
		for _, matching := range mentee_matchings {
			relationships = append(
				relationships,
				convertMatchingToRelationshipDataModel(matching.Mentee, MENTEE),
			)
		}
	}

	// find all people this person is a mentor for
	mentor_matchings, err := data.MatchingsObjs.
		Select().
		Where(data.MatchingsObjs.FilterMentee("=", c.SessionData.UserId)).
		List(c.Db)

	if err == nil && len(mentor_matchings) > 0 {
		response.State = ACCOUNT_MATCHED
		// create array of matchings
		for _, matching := range mentor_matchings {
			relationships = append(
				relationships,
				convertMatchingToRelationshipDataModel(matching.Mentor, MENTOR),
			)
		}
	}

	response.Relatationships = relationships
	c.Result = response

	return nil
}
