package meetup_reminder_test

import (
	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/ctx"
	"letstalk/server/core/meetup_reminder"
	"letstalk/server/core/sessions"
	"letstalk/server/core/test"
	"letstalk/server/core/user"
	"letstalk/server/data"

	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// Assert that one initial meetup reminder is scheduled for each user.
func assertInitialReminderScheduled(t *testing.T, db *gorm.DB, userOneId data.TUserID, userTwoId data.TUserID) {
	reminders := make([]data.MeetupReminder, 0)
	err := db.Where(&data.MeetupReminder{UserId: userOneId}).Or(&data.MeetupReminder{UserId: userTwoId}).Find(&reminders).Error
	assert.NoError(t, err)
	inOneWeekAndADay := time.Now().AddDate(0, 0, 8) // Expect reminder in ~1 week.
	assert.Len(t, reminders, 2)
	assert.Equal(t, userOneId, reminders[0].UserId)
	assert.Equal(t, userTwoId, reminders[0].MatchUserId)
	assert.Equal(t, data.MEETUP_TYPE_INITIAL, reminders[0].Type)
	assert.Equal(t, data.MEETUP_REMINDER_SCHEDULED, reminders[0].State)
	assert.True(t, reminders[0].ScheduledAt.Before(inOneWeekAndADay), "Reminder scheduled within 8 days")
	assert.Equal(t, userTwoId, reminders[1].UserId)
	assert.Equal(t, userOneId, reminders[1].MatchUserId)
	assert.Equal(t, data.MEETUP_TYPE_INITIAL, reminders[1].Type)
	assert.Equal(t, data.MEETUP_REMINDER_SCHEDULED, reminders[1].State)
	assert.True(t, reminders[1].ScheduledAt.Before(inOneWeekAndADay), "Reminder scheduled within 8 days")
}

func TestScheduleInitialReminder(t *testing.T) {
	tests := []test.Test{
		{
			TestName: "Test directly scheduling initial meetup reminders",
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				err := meetup_reminder.ScheduleInitialReminder(db, userOne.UserId, userTwo.UserId)
				assert.NoError(t, err)
				assertInitialReminderScheduled(t, db, userOne.UserId, userTwo.UserId)
			},
		},
		{
			TestName: "Test automatically scheduling reminder on new connection",
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				request := api.ConnectionRequest{
					UserId:     userTwo.UserId,
					IntentType: data.INTENT_TYPE_ASSIGNED,
				}
				connection.HandleRequestConnection(c, request)
				// Accept the request as user two.
				c.SessionData = &sessions.SessionData{UserId: userTwo.UserId}
				acceptReq := api.AcceptConnectionRequest{
					UserId: userOne.UserId,
				}
				_, err := connection.HandleAcceptConnection(c, acceptReq)
				assert.NoError(t, err)
				assertInitialReminderScheduled(t, db, userOne.UserId, userTwo.UserId)
			},
		},
	}
	test.RunTestsWithDb(tests)
}
