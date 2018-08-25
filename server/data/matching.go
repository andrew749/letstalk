package data

import (
	"github.com/jinzhu/gorm"
)

type MatchingState int

const (
	MATCHING_STATE_UNKNKOWN MatchingState = iota
	MATCHING_STATE_UNVERIFIED
	MATCHING_STATE_VERIFIED
	MATCHING_STATE_EXPIRED
)

type Matching struct {
	gorm.Model
	MentorUser *User         `gorm:"foreignkey:Mentor"`
	Mentor     TUserID       `gorm:"not null"`
	MenteeUser *User         `gorm:"foreignkey:Mentee"`
	Mentee     TUserID       `gorm:"not null"`
	State      MatchingState `gorm:"not null"`
}

func GetMatchingWithId(db *gorm.DB, matchingId uint) (*Matching, error) {
	var matching Matching
	if err := db.First(&matching, matchingId).Error; err != nil {
		return nil, err
	}
	return &matching, nil
}
