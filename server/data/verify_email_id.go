package data

import "time"

// VerifyEmailId: Generated when a user needs to verify their email address
type VerifyEmailId struct {
	Id             string    `gorm:"primary_key;unique;not null"`
	User           User      `gorm:"foreignKey:UserId"`
	UserId         TUserID   `gorm:"not null"`
	// UW email used to verify account. Not necessarily equal to the user's primary email.
	Email          string    `gorm:"type:varchar(128);not null"`
	IsActive       bool      `gorm:"not null;default=false"`
	IsUsed         bool      `gorm:"not null;default=false"`
	ExpirationDate time.Time `gorm:"not null"`
}
