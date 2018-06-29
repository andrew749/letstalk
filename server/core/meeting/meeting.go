package meeting

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/api"
	"letstalk/server/data"
	"github.com/jinzhu/gorm"
)

// PostMeetingConfirmation lets users confirm that a scheduled meeting occurred.
func PostMeetingConfirmation(c *ctx.Context) errs.Error {
	authUserId := c.SessionData.UserId
	var input api.MeetingConfirmation
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	matchedUser, err := query.GetUserBySecret(c.Db, input.Secret)
	if err != nil {
		return errs.NewRequestError("Could not find user")
	}

	// TODO: find and confirm existing meeting with this user, if exists.

	matchingObj, err := query.GetMatchingByUserIds(c.Db, authUserId, matchedUser.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if matchingObj == nil {
		return errs.NewRequestError("No existing match with this user")
	}

	// Store a confirmation of the meeting for future reference.
	conf := &data.MeetingConfirmation{
		MatchingId: matchingObj.ID,
	}

	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Model(&data.MeetingConfirmation{}).Create(conf).Error; err != nil {
			return err
		}
		// Verify the matching if this is the first confirmed meeting.
		if matchingObj.State == data.MATCHING_STATE_UNVERIFIED {
			if err := saveVerifiedMatch(tx, matchingObj); err != nil {
				return err
			}
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(err)
	}

	c.Result = input
	return nil
}

// Updates the matching to Verified state in database.
func saveVerifiedMatch(tx *gorm.DB, matching *data.Matching) error {
	matching.State = data.MATCHING_STATE_VERIFIED
	return updateMatchingObject(tx, matching)
}

// Update the matching object in database.
func updateMatchingObject(tx *gorm.DB, matching *data.Matching) error {
	if matching == nil {
		return nil
	}
	// Strip composite fields from matching struct.
	matching.MenteeUser = nil
	matching.MentorUser = nil
	return tx.Model(&data.Matching{}).UpdateColumns(matching).Error
}
