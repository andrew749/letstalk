package data

import (
	"github.com/jinzhu/gorm"
)

type UserCredential struct {
	gorm.Model
	User         User       `gorm:"foreignkey:UserId"`
	UserId       int        `gorm:"not null"`
	Credential   Credential `gorm:"foreignKey:CredentialId"`
	CredentialId uint       `gorm:"not null"`
}
