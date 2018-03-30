package data

import (
	"github.com/jinzhu/gorm"
)

type Matchings struct {
	gorm.Model
	MatchingId int  `json:"matchingId"`
	MentorUser User `gorm:"foreignkey:Mentor"`
	Mentor     int  `json:"mentor"`
	MenteeUser User `gorm:"foreignkey:Mentee"`
	Mentee     int  `json:"mentee"`
}
