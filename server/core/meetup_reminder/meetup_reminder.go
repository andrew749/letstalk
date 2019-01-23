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
	newReminder := data.MeetupReminder{
		UserId:      input.UserId,
		MatchUserId: input.MatchUserId,
		Type:        data.MEETUP_TYPE_FOLLOWUP, // POSTed reminders are always for followup.
		State:       data.MEETUP_REMINDER_SCHEDULED,
		ScheduledAt: input.ReminderTime,
	}
	dbErr := c.WithinTx(func(tx *gorm.DB) error {
		if err := tx.Delete(&data.MeetupReminder{}, &data.MeetupReminder{UserId: input.UserId, MatchUserId: input.MatchUserId}).Error; err != nil {
			return err
		}
		if err := tx.Model(&data.MeetupReminder{}).Create(&newReminder).Error; err != nil {
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
	return HandleDeleteMeetupReminder(c, input)
}

func HandleDeleteMeetupReminder(c *ctx.Context, input api.MeetupReminder) errs.Error {
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
	c.Result = "Ok"
	return nil
}

// ScheduleInitialReminder should be called whenever a new connection is made between two users. Schedules an
// initial meetup reminder for each user.
func ScheduleInitialReminder(tx *gorm.DB, userOne data.TUserID, userTwo data.TUserID) error {
	scheduledAt := time.Now().AddDate(0, 0, 7) // Schedule one week from now.
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
