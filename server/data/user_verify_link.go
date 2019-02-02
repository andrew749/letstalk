package data

import (
	"time"

	"database/sql/driver"
)

type TVerifyLinkID string

type VerifyLinkType string

func (u *VerifyLinkType) Scan(value interface{}) error {
	*u = VerifyLinkType(value.([]uint8))
	return nil
}
func (u VerifyLinkType) Value() (driver.Value, error) { return string(u), nil }

func (u *TVerifyLinkID) Scan(value interface{}) error {
	*u = TVerifyLinkID(value.([]uint8))
	return nil
}
func (u TVerifyLinkID) Value() (driver.Value, error) { return string(u), nil }

type UserVerifyLink struct {
	Id      TVerifyLinkID  `gorm:"primary_key;size:190;unique;not null"`
	User    User           `gorm:"foreignKey:UserId"`
	UserId  TUserID        `gorm:"not null"`
	Clicked bool           `gorm:"not null";default=false"`
	Type    VerifyLinkType `gorm:"not null"`

	// Times associated with link
	Times
	ExpiresAt *time.Time // optional expiry
}
