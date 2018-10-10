package data

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Direction string

const (
	MentorToMentee Direction = "MentorToMentee"
	MenteeToMentor Direction = "MenteeToMentor"
)

type SentMentorshipEmails struct {
	gorm.Model
	Mentor    string    `gorm:"not null"`
	Mentee    string    `gorm:"not null"`
	Direction string    `gorm:"not null"`
	RunId     time.Time `gorm:"not null"`
}

func SendMentorshipEmail(
	mentor string,
	mentee string,
	direction Direction,
	runId time.Time,
) error {
	return nil
}
