package matching

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/romana/rlog"
	"letstalk/server/core/sessions"
	"letstalk/server/notifications"
)

/**
 * PostMatchingController creates a new matching between two users, in an "unverified" state.
 * Only used for debugging!
 * TODO(aklen): only allow administrators to do this.
 */
func PostMatchingController(c *ctx.Context) errs.Error {
	var input api.Matching
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewClientError("Failed to parse input")
	}

	rlog.Info("Received input: ", input)
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
		State:  data.MATCHING_STATE_UNVERIFIED,
	}

	if err := c.Db.Create(matching).Error; err != nil {
		return errs.NewDbError(err)
	}

	// Send push notifications asynchronously.
	go sendMatchNotifications(c, mentor.UserId, mentee.UserId)

	c.Result = convertMatchingDataToApi(matching)
	return nil
}

// Does not populate secret field.
func convertMatchingDataToApi(matching *data.Matching) *api.Matching {
	if matching == nil {
		return nil
	}
	return &api.Matching{
		Mentor: matching.Mentor,
		Mentee: matching.Mentee,
		State:  matching.State,
	}
}

func sendMatchNotifications(
	c *ctx.Context,
	mentorId int,
	menteeId int,
) errs.Error {
	mentorDeviceTokens, err := sessions.GetDeviceTokensForUser(*c.SessionManager, mentorId)
	if err != nil {
		return errs.NewDbError(err)
	}
	menteeDeviceTokens, err := sessions.GetDeviceTokensForUser(*c.SessionManager, menteeId)
	if err != nil {
		return errs.NewDbError(err)
	}
	for _, token := range mentorDeviceTokens {
		notifications.NewMenteeNotification(token)
	}
	for _, token := range menteeDeviceTokens {
		notifications.NewMentorNotification(token)
	}
	return nil
}
