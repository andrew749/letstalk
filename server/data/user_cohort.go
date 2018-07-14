package data

type UserCohort struct {
	User     User    `gorm:"foreignkey:UserId"`
	UserId   TUserID `json:"userId" gorm:"not null;primary_key;auto_increment:false"`
	Cohort   *Cohort `gorm:"foreignkey:CohortId;association_foreignkey:CohortId"`
	CohortId int     `json:"cohortId" gorm:"not null"`
}
