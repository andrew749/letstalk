package data

import (
	"time"
)

type User struct {
	CreatedAt time.Time
	UserId    int         `json:"user_id" gorm:"not null;primary_key;auto_increment"`
	FirstName string      `json:"first_name" gorm:"not null"`
	LastName  string      `json:"last_name" gorm:"not null"`
	Email     string      `json:"email" gorm:"type:varchar(128);not null;unique"`
	Gender    int         `json:"gender" gorm:"not null"`
	Birthdate *time.Time  `json:"birthdate" gorm:"type:date;not null"`
	Sessions  []Session   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort    *UserCohort `gorm:"null"`
	Mentees   []*User     `gorm:"many2many:mentees;association_jointable_foreignkey:mentee_id"`
	Mentors   []*User     `gorm:"many2many:mentors;association_jointable_foreignkey:mentor_id"`
}
