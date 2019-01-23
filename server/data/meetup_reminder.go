package data

import (
	"database/sql/driver"
	"time"

	"github.com/jinzhu/gorm"
)

type MeetupType string

const (
	MEETUP_TYPE_INITIAL  MeetupType = "INITIAL_MEETING"
	MEETUP_TYPE_FOLLOWUP MeetupType = "FOLLOWUP_MEETING"
)

type MeetupReminderState string

const (
	MEETUP_REMINDER_SCHEDULED MeetupReminderState = "SCHEDULED"
	MEETUP_REMINDER_PROCESSED MeetupReminderState = "PROCESSED"
)

type MeetupReminder struct {
	gorm.Model
	UserId      TUserID
	MatchUserId TUserID
	Type        MeetupType
	State       MeetupReminderState
	ScheduledAt time.Time
}

func (u *MeetupType) Scan(value interface{}) error { *u = MeetupType(value.([]uint8)); return nil }
func (u MeetupType) Value() (driver.Value, error)  { return string(u), nil }

func (u *MeetupReminderState) Scan(value interface{}) error {
	*u = MeetupReminderState(value.([]uint8))
	return nil
}
func (u MeetupReminderState) Value() (driver.Value, error) { return string(u), nil }
