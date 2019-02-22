package api

import (
	"time"

	"letstalk/server/data"
)

type MeetupReminder struct {
	UserId       data.TUserID `json:"userId" binding:"required"`
	MatchUserId  data.TUserID `json:"matchUserId" binding:"required"`
	ReminderTime *time.Time   `json:"reminderTime"`
}
