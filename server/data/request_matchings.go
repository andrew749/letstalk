package data

import (
	"github.com/jinzhu/gorm"
)

type RequestMatching struct {
	gorm.Model
	AskerUser    User `gorm:"foreignkey:Asker"`
	Asker        int  `gorm:"not null"`
	AnswererUser User `gorm:"foreignkey:Answerer"`
	Answerer     int  `gorm:"not null"`
}
