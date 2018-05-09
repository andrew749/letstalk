package matching

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/api"
	"letstalk/server/data"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/utility"
	"strconv"
)

/**
 * PostMatchingController creates a new matching between two users, in an "unverified" state.
 * Only used for debugging!
 * TODO(aklen): only allow administrators to do this.
 */
func PostMatchingController(c *ctx.Context) errs.Error {
	var input api.PostMatchingRequest
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewClientError("Failed to parse input")
	}

	// Ensure both users are unique and exist.
	if input.Mentee == input.Mentor {
		return errs.NewClientError("User cannot match with themselves")
	}
	var mentor, mentee *data.User
	var err error
	if mentee, err = query.GetUserById(c.Db, input.Mentee); err != nil {
		return errs.NewClientError("Mentee not found")
	}
	if mentor, err = query.GetUserById(c.Db, input.Mentor); err != nil {
		return errs.NewClientError("Mentor not found")
	}

	// Ensure a matching doesn't already exist between these users.
	existingMatching, err := query.GetMatchingByUserIds(c.Db, mentor.UserId, mentee.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if existingMatching != nil {
		return errs.NewClientError("Matching already exists between these users")
	}

	// Ensure users have finished onboarding.
	if onboardingStatus, err := onboarding.GetOnboardingInfo(c.Db, mentor.UserId); err != nil {
		return err
	} else if onboardingStatus.State != api.ONBOARDING_DONE {
		return errs.NewClientError("Mentor is not finished onboarding")
	}
	if onboardingStatus, err := onboarding.GetOnboardingInfo(c.Db, mentee.UserId); err != nil {
		return err
	} else if onboardingStatus.State != api.ONBOARDING_DONE {
		return errs.NewClientError("Mentee is not finished onboarding")
	}

	// Insert new matching.
	matching := &data.Matching{
		Mentee: mentee.UserId,
		Mentor: mentor.UserId,
		State: api.MATCHING_STATE_UNVERIFIED,
	}

	if matching.MenteeSecret, err = getNewMatchingSecret(); err != nil {
		return errs.NewInternalError("Error generating matching secret: %v", err)
	}
	if matching.MentorSecret, err = getNewMatchingSecret(); err != nil {
		return errs.NewInternalError("Error generating matching secret: %v", err)
	}

	if err := c.Db.Create(matching).Error; err != nil {
		return errs.NewDbError(err)
	}

	result := api.MatchingResult{
		Mentee: matching.Mentee,
		Mentor: matching.Mentor,
		State: matching.State,
	}
	c.Result = &result

	return nil
}

func getNewMatchingSecret() (string, error) {
	return utility.GenerateRandomString(20)
}

// GetMatchingController gets details for a match with the authenticated user.
func GetMatchingController(c *ctx.Context) errs.Error {
	inputUserId := c.GinContext.Param("user_id")
	if len(inputUserId) == 0 {
		return errs.NewClientError("No user id given")
	}
	matchUserId, err := strconv.Atoi(inputUserId)
	if err != nil {
		return errs.NewClientError("User id in unexpected format")
	}
	authUserId := c.SessionData.UserId
	matchingObj, err := query.GetMatchingByUserIds(c.Db, authUserId, matchUserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	result := &api.MatchingResult{
		Mentor: matchingObj.Mentor,
		Mentee: matchingObj.Mentee,
		State: matchingObj.State,
	}
	if matchingObj.Mentor == authUserId {
		// Auth user is the mentor.
		result.Secret = matchingObj.MentorSecret
	} else {
		// Auth user is the mentee.
		result.Secret = matchingObj.MenteeSecret
	}
	c.Result = result
	return nil
}
