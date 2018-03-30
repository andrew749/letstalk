package data

import (
	"time"
)

type FbAuthToken struct {
	User      User      `gorm:"foreignkey:UserId"`
	UserId    int       `json:"userId" gorm:"not null;primary_key"`
	AuthToken string    `json:"authToken" gorm:"not null"`
	Expiry    time.Time `json:"expiry" gorm:"not null"`
}
