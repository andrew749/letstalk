package data

import (
	"encoding/json"
	"time"

	"database/sql/driver"

	"github.com/jinzhu/gorm"
)

type NotifType string

type NotifState string

const (
	NOTIF_STATE_UNREAD NotifState = "UNREAD"
	NOTIF_STATE_READ              = "READ"
)

type JSONBlob json.RawMessage

type Notification struct {
	gorm.Model
	UserId        TUserID    `gorm:"not null"`
	User          User       `gorm:"foreignkey:UserId"`
	Type          NotifType  `gorm:"not null"`
	Timestamp     time.Time  `gorm:"not null;default:now()"` // when the notification was created in the system (not in db)
	State         NotifState `gorm:"not null;default:PENDING_SEND"`
	Title         string     `gorm:"not null"`
	Message       string     `gorm:"not null"`
	ThumbnailLink *string    `gorm:""`
	Data          JSONBlob   `gorm:"not null" sql:"type:json"`
}

func (u *NotifType) Scan(value interface{}) error { *u = NotifType(value.([]byte)); return nil }
func (u NotifType) Value() (driver.Value, error)  { return string(u), nil }

func (u *NotifState) Scan(value interface{}) error { *u = NotifState(value.([]byte)); return nil }
func (u NotifState) Value() (driver.Value, error)  { return string(u), nil }

func (u *JSONBlob) Scan(value interface{}) error {
	*u = JSONBlob(value.([]byte))
	return nil
}
func (u JSONBlob) Value() (driver.Value, error) { return string(u), nil }

// PendingNotifications Notifications that have been sent to expo but not necessarily delivered
type ExpoPendingNotification struct {
	gorm.Model
	Notification   Notification `gorm:"foreign_key:NotificationId"`
	NotificationId uint         `gorm:"primary_key;auto_increment:false"`
	DeviceId       string       `gorm:"not null;primary_key;"`
	Receipt        *string      `gorm:""`
	FailureMessage *string
	FailureDetails *string
}

func CreateNewPendingNotification(db *gorm.DB, notificationId uint, deviceId string) (*ExpoPendingNotification, error) {
	notification := ExpoPendingNotification{
		NotificationId: notificationId,
		DeviceId:       deviceId,
	}

	if err := db.Create(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (e *ExpoPendingNotification) MarkNotificationError(db *gorm.DB, errorMessage *string, errorDetails interface{}) error {
	serializedErrorDetails, err := json.Marshal(errorDetails)
	if err != nil {
		return err
	}
	serializedErrorString := string(serializedErrorDetails)
	e.FailureDetails = &serializedErrorString
	e.FailureMessage = errorMessage
	return db.Save(e).Error
}

func (e *ExpoPendingNotification) MarkNotificationSent(db *gorm.DB, receipt string) error {
	e.Receipt = &receipt
	return db.Save(e).Error
}
