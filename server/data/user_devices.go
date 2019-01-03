package data

import (
	"database/sql/driver"

	"github.com/jinzhu/gorm"
)

type NotificationTokenType string

func (u *NotificationTokenType) Scan(value interface{}) error {
	*u = NotificationTokenType(value.([]byte))
	return nil
}
func (u NotificationTokenType) Value() (driver.Value, error) { return string(u), nil }

const (
	EXPO_PUSH = "EXPO"
)

type UserDevice struct {
	User                  User                  `gorm:"foreign_key:UserId"`
	UserId                TUserID               `gorm:"primary_key"`
	NotificationToken     string                `gorm:"size:190;primary_key;not null"`
	NotificationTokenType NotificationTokenType `gorm:"not null"`
}

func AddDeviceTokenForUser(db *gorm.DB, userId TUserID, token string, tokenType NotificationTokenType) error {
	userDevice := UserDevice{
		UserId:                userId,
		NotificationToken:     token,
		NotificationTokenType: tokenType,
	}

	return db.FirstOrCreate(&userDevice).Error
}

func AddExpoDeviceTokenForUser(db *gorm.DB, userId TUserID, token string) error {
	return AddDeviceTokenForUser(db, userId, token, EXPO_PUSH)
}

func GetDevicesForUser(db *gorm.DB, userId TUserID) (*[]UserDevice, error) {
	var userDevices []UserDevice
	if err := db.Where("user_id=?", userId).Find(&userDevices).Error; err != nil {
		return nil, err
	}
	return &userDevices, nil
}

func GetDeviceNotificationTokensForUser(db *gorm.DB, userId TUserID) (*[]string, error) {
	devices, err := GetDevicesForUser(db, userId)
	if err != nil {
		return nil, err
	}

	tokens := make([]string, len(*devices))
	for i, device := range *devices {
		tokens[i] = device.NotificationToken
	}

	return &tokens, nil
}
