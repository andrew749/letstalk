package data

import (
	"github.com/jinzhu/gorm"
)

type UserVector struct {
	gorm.Model
	UserId         int  `json:"userId" gorm:"not null"`
	User           User `gorm:"foreignkey:UserId"`
	PreferenceType int  `json:"preferenceType" gorm:"not null"`
	Sociable       int  `json:"sociable" gorm:"not null"`
	HardWorking    int  `json:"hardworking" gorm:"not null"`
	Ambitious      int  `json:"ambitious" gorm:"not null"`
	Energetic      int  `json:"energetic" gorm:"not null"`
	Carefree       int  `json:"carefree" gorm:"not null"`
	Confident      int  `json:"confident" gorm:"not null"`
}
