package data

import (
	"github.com/jinzhu/gorm"
)

type UserCredentialRequest struct {
	gorm.Model
	User         User       `gorm:"foreignkey:UserId"`
	UserId       int        `json:"userId" gorm:"not null"`
	Credential   Credential `gorm:"foreignKey:CredentialId"`
	CredentialId uint       `gorm:"not null"`
}
