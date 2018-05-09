package data

import (
	"github.com/jinzhu/gorm"
)

type MatchingState int

const (
	MATCHING_STATE_UNVERIFIED MatchingState = iota
	MATCHING_STATE_VERIFIED
	MATCHING_STATE_EXPIRED
)

type Matching struct {
	gorm.Model
	MentorUser *User `gorm:"foreignkey:Mentor"`
	Mentor     int  `gorm:"not null"`
	MentorSecret string  `gorm:"not null"`
	MenteeUser *User `gorm:"foreignkey:Mentee"`
	Mentee     int  `gorm:"not null"`
	MenteeSecret string  `gorm:"not null"`
	State MatchingState `gorm:"not null"`
}
