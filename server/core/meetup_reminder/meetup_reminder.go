package meetup_reminder

import (
	"time"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

const NUM_DAYS_UNTIL_INITIAL_REMINDER = 7 // Schedule initial meetup reminder a week after match.

// PostMeetupReminder replaces existing meetup reminders for an ordered (user, match) pair with the given reminder.
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
	return HandlePostMeetupReminder(c, input)
}

func HandlePostMeetupReminder(c *ctx.Context, input api.MeetupReminder) errs.Error {
	if input.ReminderTime == nil {
		return errs.NewRequestError("Missing reminder time field")
	}
	newReminder := data.MeetupReminder{
		UserId:      input.UserId,
		MatchUserId: input.MatchUserId,
		Type:        data.MEETUP_TYPE_FOLLOWUP, // POSTed reminders are always for followup.
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: *input.ReminderTime,
	}
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Model(&data.MeetupReminder{}).Where(&data.MeetupReminder{UserId: input.UserId, MatchUserId: input.MatchUserId, State: data.MEETUP_REMINDER_SCHEDULED}).
		Update(&data.MeetupReminder{State: data.MEETUP_REMINDER_REPLACED}).Error; err != nil {
			return err
		}
		if err := tx.Create(&newReminder).Error; err != nil {
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

// DeleteMeetupReminder cancels existing meetup reminders for (user, match) and (match, user).
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
	return HandleCancelMeetupReminder(c, input)
}

func HandleCancelMeetupReminder(c *ctx.Context, input api.MeetupReminder) errs.Error {
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Model(&data.MeetupReminder{}).Where(&data.MeetupReminder{UserId: input.UserId, MatchUserId: input.MatchUserId, State: data.MEETUP_REMINDER_SCHEDULED}).
		Update(&data.MeetupReminder{State: data.MEETUP_REMINDER_CANCELLED}).Error; err != nil {
			return err
		}
		if err := tx.Model(&data.MeetupReminder{}).Where(&data.MeetupReminder{UserId: input.MatchUserId, MatchUserId: input.UserId, State: data.MEETUP_REMINDER_SCHEDULED}).
		Update(&data.MeetupReminder{State: data.MEETUP_REMINDER_CANCELLED}).Error; err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	c.Result = "Ok"
	return nil
}

// ScheduleInitialReminder should be called whenever a new connection is made between two users. Schedules an
// initial meetup reminder for each user.
func ScheduleInitialReminder(tx *gorm.DB, userOne data.TUserID, userTwo data.TUserID) error {
	scheduledAt := time.Now().AddDate(0, 0, NUM_DAYS_UNTIL_INITIAL_REMINDER)
	userOneReminder := data.MeetupReminder{
		UserId:      userOne,
		MatchUserId: userTwo,
		Type:        data.MEETUP_TYPE_INITIAL,
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: scheduledAt,
	}
	userTwoReminder := data.MeetupReminder{
		UserId:      userTwo,
		MatchUserId: userOne,
		Type:        data.MEETUP_TYPE_INITIAL,
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: scheduledAt,
	}
	if err := tx.Create(&userOneReminder).Error; err != nil {
		return err
	}
	if err := tx.Create(&userTwoReminder).Error; err != nil {
		return err
	}
	return nil
}
