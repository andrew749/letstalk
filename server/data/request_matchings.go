package data

import (
	"github.com/jinzhu/gorm"
)

type RequestMatching struct {
	gorm.Model
	AskerUser    *User       `gorm:"foreignkey:Asker"`
	Asker        TUserID     `gorm:"not null"`
	AnswererUser *User       `gorm:"foreignkey:Answerer"`
	Answerer     TUserID     `gorm:"not null"`
	Credential   *Credential `gorm:"foreignkey:CredentialId"`
	CredentialId uint        `gorm:"not null"`
}
