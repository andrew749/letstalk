package data

import (
	"github.com/jinzhu/gorm"
)

type NotificationTokenType string

const (
	EXPO_PUSH = "expo"
)

type UserDevice struct {
	User                  User    `gorm:"foreign_key:UserId"`
	UserId                TUserID `gorm:"primary_key"`
	NotificationToken     string  `gorm:"not null"`
	NotificationTokenType NotificationTokenType
}

func AddDeviceTokenForUser(db *gorm.DB, userId TUserID, token string, tokenType NotificationTokenType) error {
	userDevice := UserDevice{
		UserId:                userId,
		NotificationToken:     token,
		NotificationTokenType: tokenType,
	}

	return db.FirstOrCreate(&userDevice).Error
}

func AddExpoDeviceTokenforUser(db *gorm.DB, userId TUserID, token string) error {
	return AddDeviceTokenForUser(db, userId, token, EXPO_PUSH)
}

func GetDeviceNotificationTokensForUser(db *gorm.DB, userId TUserID) (*[]string, error) {
	var users []string
	if err := db.Where("user_id=?", userId).Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
