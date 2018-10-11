package data

import (
	"database/sql/driver"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type Direction string

const (
	ToMentor Direction = "ToMentor"
	ToMentee Direction = "ToMentee"
)

func (u *Direction) Scan(value interface{}) error { *u = Direction(value.([]byte)); return nil }
func (u Direction) Value() (driver.Value, error)  { return string(u), nil }

type SentMentorshipEmails struct {
	gorm.Model
	MentorID   TUserID   `gorm:"not null"`
	MentorName string    `gorm:"not null"`
	MenteeID   TUserID   `gorm:"not null"`
	MenteeName string    `gorm:"not null"`
	Direction  Direction `gorm:"not null"`
	RunId      time.Time `gorm:"not null"`
}

func SendMentorshipEmail(
	db *gorm.DB,
	mentorId TUserID,
	mentorName string,
	menteeId TUserID,
	menteeName string,
	direction Direction,
	runId time.Time,
) error {
	sentEmail := SentMentorshipEmails{
		MentorID:   mentorId,
		MentorName: mentorName,
		MenteeID:   menteeId,
		MenteeName: menteeName,
		Direction:  direction,
		RunId:      runId,
	}
	if err := db.Create(&sentEmail).Error; err != nil {
		rlog.Criticalf("Unable to create new email sent Record: %+v", err)
		return err
	}
	return nil
}
