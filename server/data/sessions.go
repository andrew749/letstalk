package data

import (
	"time"
)

type Session struct {
	SessionId         string    `json:"sessionId"`
	User              User      `gorm:"foreignkey:UserId"`
	UserId            int       `json:"userId"`
	ExpiryDate        time.Time `json:"expiryDate"`
	NotificationToken *NotificationToken
}
