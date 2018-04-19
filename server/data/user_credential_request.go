package data

import (
	"github.com/jinzhu/gorm"
)

type UserCredentialRequest struct {
	gorm.Model
	User           User                     `gorm:"foreignkey:UserId"`
	UserId         int                      `json:"userId" gorm:"not null"`
	PositionId     CredentialPositionId     `json:"positionId" gorm:"not null"`
	OrganizationId CredentialOrganizationId `json:"organizationId" gorm:"not null"`
}
