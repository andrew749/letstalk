package data

import (
	"github.com/jinzhu/gorm"
)

type UserCredential struct {
	gorm.Model
	User           User `gorm:"foreignkey:UserId"`
	UserId         int  `json:"userId" gorm:"not null"`
	PositionId     int  `json:"positionId" gorm:"not null"`
	OrganizationId int  `json:"organizationId" gorm:"not null"`
}
