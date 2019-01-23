package query

import (
	"time"

	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetMeetupRemindersScheduledBefore(db *gorm.DB, before time.Time) ([]data.MeetupReminder, error) {
	var reminders []data.MeetupReminder
	q := db.Model(&data.MeetupReminder{}).
			Where(&data.MeetupReminder{State: data.MEETUP_REMINDER_SCHEDULED}).
			Where("meetup_reminders.scheduled_at <= ?", before)
	if err := q.Find(&reminders).Error; err != nil {
		return nil, err
	}
	return reminders, nil
}
