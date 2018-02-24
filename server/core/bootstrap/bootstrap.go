package bootstrap

import (
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

type BootstrapResponse struct {
	State BootstrapState `json:"state"`
}

/**
 * Returns what the current status of a user is
 */
func GetCurrentUserBoostrapStatusController(c *ctx.Context) errs.Error {
	// since this method is authenticated the account needs to exist.
	var response = BootstrapResponse{
		State: ACCOUNT_CREATED,
	}

	// check if the user has been matched with another user yet
	_, err := data.MatchingsObjs.
		Select().
		Where(data.MatchingsObjs.FilterUser("=", c.SessionData.UserId)).
		List(c.Db)

	if err == nil {
		response.State = ACCOUNT_MATCHED
		c.Result = response
		return nil
	}

	// check if the account has been onboarded
	_, err = data.UserCohortObjs.
		Select().
		Where(data.UserCohortObjs.FilterUserId("=", c.SessionData.UserId)).
		One(c.Db)

	if err == nil {
		response.State = ACCOUNT_SETUP
		c.Result = response
		return nil
	}

	return nil
}
