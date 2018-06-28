package data

import (
	"github.com/jinzhu/gorm"
)

type MeetingConfirmation struct {
	gorm.Model
	Matching   *Matching `gorm:"foreignkey:MatchingId"`
	MatchingId uint `gorm:"not null"`
}
