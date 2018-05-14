package meeting

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/api"
	"letstalk/server/data"
)

// PostMeetingConfirmation lets users confirm that a scheduled meeting occurred.
func PostMeetingConfirmation(c *ctx.Context) errs.Error {
	authUserId := c.SessionData.UserId
	var input api.MeetingConfirmation
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewClientError("Failed to parse input")
	}
	matchedUser, err := query.GetUserBySecret(c.Db, input.Secret)
	if err != nil {
		return errs.NewClientError("Could not find user")
	}

	// TODO: find and confirm existing meeting with this user, if exists.

	matchingObj, err := query.GetMatchingByUserIds(c.Db, authUserId, matchedUser.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if matchingObj == nil {
		return errs.NewClientError("No existing match with this user")
	}

	// Verify the matching if this is the first confirmed meeting.
	if matchingObj.State == data.MATCHING_STATE_UNVERIFIED {
		// TODO: do in transaction
		if err := saveVerifiedMatch(c, matchingObj); err != nil {
			return err
		}
	}

	c.Result = input
	return nil
}

// Updates the matching to Verified state in database.
func saveVerifiedMatch(c *ctx.Context, matching *data.Matching) errs.Error {
	matching.State = data.MATCHING_STATE_VERIFIED
	if err := updateMatchingObject(c, matching); err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

// Update the matching object in database.
func updateMatchingObject(c *ctx.Context, matching *data.Matching) error {
	if matching == nil {
		return nil
	}
	// Strip composite fields from matching struct.
	matching.MenteeUser = nil
	matching.MentorUser = nil
	return c.Db.Model(&data.Matching{}).UpdateColumns(matching).Error
}
