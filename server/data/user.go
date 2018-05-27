package data

import (
	"time"
)

type User struct {
	CreatedAt        time.Time           `gorm:"not null"`
	UserId           int                 `gorm:"not null;primary_key;auto_increment"`
	FirstName        string              `gorm:"not null"`
	LastName         string              `gorm:"not null"`
	Email            string              `gorm:"type:varchar(128);not null;unique"`
	Secret           string              `gorm:"type:char(36);not null;unique"`
	Gender           int                 `gorm:"not null"`
	Birthdate        *time.Time          `gorm:"type:date;not null"`
	Sessions         []Session           `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AuthData         *AuthenticationData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	ExternalAuthData *ExternalAuthData   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort           *UserCohort         `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AdditionalData   *UserAdditionalData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}
