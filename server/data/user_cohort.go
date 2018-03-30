package data

import (
	"github.com/jinzhu/gorm"
)

type UserCohort struct {
	gorm.Model
	User     User   `gorm:"foreignkey:UserId"`
	UserId   int    `json:"userId"`
	Cohort   Cohort `gorm:"foreignkey:CohortId"`
	CohortId int    `json:"cohortId"`
}
