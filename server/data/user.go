package data

import (
	"time"
	"database/sql/driver"
)

type UserRole string

const (
	USER_ROLE_DEFAULT UserRole = "DEFAULT"
	USER_ROLE_ADMIN   UserRole = "ADMIN"
)

type User struct {
	CreatedAt        time.Time `gorm:"not null"`
	UserId           int       `gorm:"not null;primary_key;auto_increment"`
	FirstName        string    `gorm:"not null"`
	LastName         string    `gorm:"not null"`
	Email            string    `gorm:"type:varchar(128);not null;unique"`
	Secret           string    `gorm:"type:char(36);not null;unique"`
	Gender           int       `gorm:"not null"`
	Birthdate        string    `gorm:"type:varchar(100);not null"`
	Role             UserRole  `gorm:"not null"`
	ProfilePic       *string
	Sessions         []Session           `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AuthData         *AuthenticationData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	ExternalAuthData *ExternalAuthData   `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	Cohort           *UserCohort         `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
	AdditionalData   *UserAdditionalData `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}

func (u *UserRole) Scan(value interface{}) error { *u = UserRole(value.([]byte)); return nil }
func (u UserRole) Value() (driver.Value, error)  { return string(u), nil }
