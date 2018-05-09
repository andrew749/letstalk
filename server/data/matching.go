package data

import (
	"github.com/jinzhu/gorm"
	"letstalk/server/core/api"
)

type Matching struct {
	gorm.Model
	MentorUser User `gorm:"foreignkey:Mentor"`
	Mentor     int  `gorm:"not null"`
	MentorSecret string  `gorm:"not null"`
	MenteeUser User `gorm:"foreignkey:Mentee"`
	Mentee     int  `gorm:"not null"`
	MenteeSecret string  `gorm:"not null"`
	State api.MatchingState `gorm:"not null"`
}
