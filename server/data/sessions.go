package data

import (
	"time"
)

type Session struct {
	SessionId         string    `json:"session_id"`
	User              User      `gorm:"foreignkey:UserId"`
	UserId            int       `json:"user_id"`
	ExpiryDate        time.Time `json:"expiry_date"`
	NotificationToken *NotificationToken
}
