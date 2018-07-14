package data

import (
	"time"
)

type Session struct {
	SessionId         string             `gorm:"not null;primary_key"`
	User              User               `gorm:"foreignkey:UserId"`
	UserId            TUserID            `gorm:"not null"`
	ExpiryDate        time.Time          `gorm:"not null"`
	NotificationToken *NotificationToken `gorm:"foreignkey:SessionId;association_foreignkey:SessionId"`
}
