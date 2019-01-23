package meeting

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"letstalk/server/core/connection"
)

// PostMeetupReminder replaces exisitng meetup reminders for an ordered (user, match) pair with the given reminder.
func PostMeetupReminder(c *ctx.Context) errs.Error {
	authUser, err := query.GetUserById(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	var input api.MeetupReminder
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if authUser.UserId != input.UserId {
		return errs.NewUnauthorizedError("Not authorized")
	}

	newReminder := data.MeetupReminder{
		UserId:      input.UserId,
		MatchUserId: input.MatchUserId,
		Type:        data.MEETUP_TYPE_FOLLOWUP, // POSTed reminders are always for followup.
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: input.ReminderTime,
	}
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Delete(data.MeetupReminder{}, data.MeetupReminder{UserId: input.UserId, MatchUserId: input.MatchUserId}).Error; err != nil {
			return err
		}
		if err := tx.Model(data.MeetupReminder{}).Create(newReminder).Error; err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}

	c.Result = input

	return nil
}

// DeleteMeetupReminder deletes existing meetup reminders for (user, match) and (match, user).
func DeleteMeetupReminder(c *ctx.Context) errs.Error {
	authUser, err := query.GetUserById(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	var input api.MeetupReminder
	if err := c.GinContext.BindJSON(&input); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}
	if authUser.UserId != input.UserId {
		return errs.NewUnauthorizedError("Not authorized")
	}

	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Delete(data.MeetupReminder{}, data.MeetupReminder{UserId: input.UserId, MatchUserId: input.MatchUserId}).Error; err != nil {
			return err
		}
		if err := tx.Delete(data.MeetupReminder{}, data.MeetupReminder{UserId: input.MatchUserId, MatchUserId: input.UserId}).Error; err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}

	c.Result = input
	return nil
}

