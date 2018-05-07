package data

import (
	"github.com/jinzhu/gorm"
)

type Matching struct {
	gorm.Model
	MentorUser User `gorm:"foreignkey:Mentor"`
	Mentor     int  `gorm:"not null"`
	MenteeUser User `gorm:"foreignkey:Mentee"`
	Mentee     int  `gorm:"not null"`
}
