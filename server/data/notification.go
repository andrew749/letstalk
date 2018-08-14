package data

import (
	"encoding/json"

	"database/sql/driver"
	"github.com/jinzhu/gorm"
)

type NotifType string

const (
	NOTIF_TYPE_NEW_CREDENTIAL_MATCH NotifType = "NEW_CREDENTIAL_MATCH"
)

type NotifState string

const (
	NOTIF_STATE_UNREAD NotifState = "UNREAD"
	NOTIF_STATE_READ              = "READ"
)

type JSONBlob json.RawMessage

type Notification struct {
	gorm.Model
	UserId TUserID    `gorm:"not null"`
	User   User       `gorm:"foreignkey:UserId"`
	Type   NotifType  `gorm:"not null"`
	State  NotifState `gorm:"not null"`
	Data   JSONBlob   `gorm:"not null" sql:"type:json"`
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
