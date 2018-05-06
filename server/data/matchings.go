package data

import (
	"github.com/jinzhu/gorm"
)

type Matchings struct {
	gorm.Model
	MentorUser User `gorm:"foreignkey:Mentor"`
	Mentor     int  `json:"mentor" gorm:"not null"`
	MenteeUser User `gorm:"foreignkey:Mentee"`
	Mentee     int  `json:"mentee" gorm:"not null"`
}
