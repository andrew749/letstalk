package meeting

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/api"
	"letstalk/server/data"
	"github.com/jinzhu/gorm"
	"letstalk/server/core/sessions"
	"letstalk/server/notifications"
	"github.com/romana/rlog"
	"github.com/getsentry/raven-go"
)

// PostMeetingConfirmation lets users confirm that a scheduled meeting occurred.
func PostMeetingConfirmation(c *ctx.Context) errs.Error {
	authUser, err := query.GetUserById(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	var input api.MeetingConfirmation
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	matchedUser, err := query.GetUserBySecret(c.Db, input.Secret)
	if err != nil {
		return errs.NewRequestError("Could not find user")
	}

	// TODO: find and confirm existing meeting with this user, if exists.

	matchingObj, err := query.GetMatchingByUserIds(c.Db, authUser.UserId, matchedUser.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}
	if matchingObj == nil {
		return errs.NewRequestError("No existing match with this user")
	}
	isFirstMeeting := matchingObj.State == data.MATCHING_STATE_UNVERIFIED

	// Store a confirmation of the meeting for future reference.
	conf := &data.MeetingConfirmation{
		MatchingId: matchingObj.ID,
	}

	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Model(&data.MeetingConfirmation{}).Create(conf).Error; err != nil {
			return err
		}
		// Verify the matching if this is the first confirmed meeting.
		if isFirstMeeting {
			if err := saveVerifiedMatch(tx, *matchingObj); err != nil {
				return err
			}
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(err)
	}

	if isFirstMeeting {
		// Also send a notification now that the match is verified.
		go func() {
			if err := sendMatchVerifiedNotifications(c, authUser, matchedUser); err != nil {
				rlog.Errorf("Error sending verified match notification: %s", err)
				raven.CaptureError(err, nil)
			}
		}()
	}

	c.Result = input
	return nil
}

// Updates the matching to Verified state in database.
func saveVerifiedMatch(tx *gorm.DB, matching data.Matching) error {
	matching.State = data.MATCHING_STATE_VERIFIED
	return updateMatchingObject(tx, matching)
}

// Update the matching object in database.
func updateMatchingObject(tx *gorm.DB, matching data.Matching) error {
	// Strip composite fields from matching struct.
	matching.MenteeUser = nil
	matching.MentorUser = nil
	return tx.Model(&data.Matching{}).UpdateColumns(matching).Error
}

// Send notifications to the two users in a newly verified match.
func sendMatchVerifiedNotifications(c *ctx.Context, verifyingUser *data.User, matchedUser *data.User) error {
	verifierDeviceTokens, err := sessions.GetDeviceTokensForUser(*c.SessionManager, verifyingUser.UserId)
	if err != nil {
		return err
	}
	matchedDeviceTokens, err := sessions.GetDeviceTokensForUser(*c.SessionManager, matchedUser.UserId)
	if err != nil {
		return err
	}
	for _, token := range verifierDeviceTokens {
		if err := notifications.MatchVerifiedNotification(token, matchedUser.FirstName); err != nil {
			return err
		}
	}
	for _, token := range matchedDeviceTokens {
		if err := notifications.MatchVerifiedNotification(token, verifyingUser.FirstName); err != nil {
			return err
		}
	}
	return nil
}
