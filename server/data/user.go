package data

import (
	"time"
)

type User struct {
	CreatedAt time.Time   `json:"createdAt" gorm:"not null"`
	UserId    int         `json:"userId" gorm:"not null;primary_key;auto_increment"`
	FirstName string      `json:"firstName" gorm:"not null"`
	LastName  string      `json:"lastName" gorm:"not null"`
	Email     string      `json:"email" gorm:"type:varchar(128);not null;unique"`
	Gender    int         `json:"gender" gorm:"not null"`
	Birthdate *time.Time  `json:"birthdate" gorm:"type:date;not null"`
	Sessions  []Session   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort    *UserCohort `gorm:"null"`
	Answerers []*User     `gorm:"many2many:answerers;association_jointable_foreignkey:answerer_id"`
	Askers    []*User     `gorm:"many2many:askers;association_jointable_foreignkey:asker_id"`
}
