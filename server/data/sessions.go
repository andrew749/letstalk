package data

import (
	"time"
)

type Session struct {
	SessionId  string    `json:"sessionId" gorm:"not null;primary_key"`
	User       User      `gorm:"foreignkey:UserId"`
	UserId     int       `json:"userId" gorm:"not null"`
	ExpiryDate time.Time `json:"expiryDate" gorm:"not null"`
	// TODO: Wtf. fix this.
	NotificationToken *NotificationToken
}
