package data

import (
	"time"
)

type FbAuthToken struct {
	User      User      `gorm:"foreignkey:UserId"`
	UserId    int       `json:"userId" gorm:"primary_key"`
	AuthToken string    `json:"authToken"`
	Expiry    time.Time `json:"expiry"`
}
