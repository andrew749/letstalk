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
	NOTIF_STATE_UNREAD       NotifState = "UNREAD"
	NOTIF_STATE_READ                    = "READ"
	NOTIF_STATE_PENDING_SEND            = "PENDING_SEND"
)

type JSONBlob json.RawMessage

type Notification struct {
	gorm.Model
	UserId        TUserID    `gorm:"not null"`
	User          User       `gorm:"foreignkey:UserId"`
	Type          NotifType  `gorm:"not null"`
	State         NotifState `gorm:"not null"`
	Timestamp     time.Time  `gorm:"not null;default:now()"` // when the notification was created in the system (not in db)
	Title         string     `gorm:"not null"`
	Message       string     `gorm:"not null"`
	ThumbnailLink *string    `gorm:""`
	Data          JSONBlob   `gorm:"not null" sql:"type:json"`
	Receipt       *string    `gorm:""`
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
