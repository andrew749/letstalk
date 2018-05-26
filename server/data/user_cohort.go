package data

type UserCohort struct {
	User     User    `gorm:"foreignkey:UserId"`
	UserId   int     `json:"userId" gorm:"not null;primary_key;auto_increment:false"`
	Cohort   *Cohort `gorm:"foreignkey:CohortId"`
	CohortId int     `json:"cohortId" gorm:"not null"`
}
