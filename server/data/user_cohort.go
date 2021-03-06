package data

type BelongingToUser interface {
	GetUser() *User
}

// TODO: Maybe denormalize cohort data here so that we don't need two table lookups to get
// cohort.
type UserCohort struct {
	UserId   TUserID   `gorm:"not null;primary_key;auto_increment:false"`
	CohortId TCohortID `gorm:"not null"`
	Times

	User   *User   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort *Cohort `gorm:"foreignkey:CohortId;association_foreignkey:CohortId"`
}

func (cohort *UserCohort) GetUser() *User {
	return cohort.User
}
