package data

import (
	"database/sql/driver"
)

type (
	TUserGroupID EntID
	TGroupID     string
)

func (u *TGroupID) Scan(value interface{}) error { *u = TGroupID(value.([]byte)); return nil }
func (u TGroupID) Value() (driver.Value, error) { return string(u), nil }

// Stores groups that a user is a part of
type UserGroup struct {
	Id        TUserGroupID `gorm:"primary_key;not null;auto_increment:true"`
	UserId    TUserID      `gorm:"not null"`
	GroupId   TGroupID     `gorm:"not null"`
	GroupName string       `gorm:"not null"`
	Times

	User *User `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}
