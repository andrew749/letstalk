package data

import (
	"github.com/jinzhu/gorm"
)

type UserCohort struct {
	gorm.Model
	User     User   `gorm:"foreignkey:UserId"`
	UserId   int    `json:"user_id"`
	Cohort   Cohort `gorm:"foreignkey:CohortId"`
	CohortId int    `json:"cohort_id"`
}
