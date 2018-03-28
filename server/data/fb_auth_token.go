package data

import (
	"time"
)

type FbAuthToken struct {
	User      User      `gorm:"foreignkey:UserId"`
	UserId    int       `json:"user_id" gorm:"primary_key"`
	AuthToken string    `json:"auth_token"`
	Expiry    time.Time `json:"expiry"`
}
