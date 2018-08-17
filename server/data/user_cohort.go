package data

// TODO: Maybe denormalize cohort data here so that we don't need two table lookups to get
// cohort.
type UserCohort struct {
	User     User      `gorm:"foreignkey:UserId"`
	UserId   TUserID   `gorm:"not null;primary_key;auto_increment:false"`
	CohortId TCohortID `gorm:"not null"`
	Times
	Cohort *Cohort `gorm:"foreignkey:CohortId;association_foreignkey:CohortId"`
}
