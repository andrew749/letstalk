package matching

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/core/onboarding"
	"letstalk/server/core/query"
	"letstalk/server/data"

	raven "github.com/getsentry/raven-go"
	"github.com/romana/rlog"
)

/**
 * PostMatchingController creates a new matching between two users, in an "unverified" state.
 * Only used for debugging!
 * TODO(aklen): only allow administrators to do this.
 */
func PostMatchingController(c *ctx.Context) errs.Error {
	var input api.Matching
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	rlog.Info("Received input: ", input)
	// Ensure both users are unique and exist.
	if input.Mentee == input.Mentor {
		return errs.NewRequestError("User cannot match with themselves")
	}
	var mentor, mentee *data.User
	var err error
	if mentee, err = query.GetUserById(c.Db, input.Mentee); err != nil {
		return errs.NewNotFoundError("Mentee not found")
	}
	if mentor, err = query.GetUserById(c.Db, input.Mentor); err != nil {
		return errs.NewNotFoundError("Mentor not found")
	}

	// Ensure a matching doesn't already exist between these users.
	existingMatching, err := query.GetMatchingByUserIds(c.Db, mentor.UserId, mentee.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if existingMatching != nil {
		return errs.NewRequestError("Matching already exists between these users")
	}

	// Ensure users have finished onboarding.
	if onboardingStatus, err := onboarding.GetOnboardingInfo(c.Db, mentor.UserId); err != nil {
		return err
	} else if onboardingStatus.State != api.ONBOARDING_DONE {
		return errs.NewRequestError("Mentor is not finished onboarding")
	}
	if onboardingStatus, err := onboarding.GetOnboardingInfo(c.Db, mentee.UserId); err != nil {
		return err
	} else if onboardingStatus.State != api.ONBOARDING_DONE {
		return errs.NewRequestError("Mentee is not finished onboarding")
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
	mentorId data.TUserID,
	menteeId data.TUserID,
) errs.Error {
	err1 := notifications.NewMenteeNotification(c.Db, menteeId)
	err2 := notifications.NewMentorNotification(c.Db, mentorId)
	if err1 != nil {
		raven.CaptureError(err1, nil)
	}
	if err2 != nil {
		raven.CaptureError(err2, nil)
	}
	return nil
}
