package data

import (
	"database/sql/driver"
)

type (
	TUserGroupID EntID
	TGroupID     string
)

func (u *TGroupID) Scan(value interface{}) error { *u = TGroupID(value.([]byte)); return nil }
func (u TGroupID) Value() (driver.Value, error)  { return string(u), nil }

// Stores groups that a user is a part of
type UserGroup struct {
	Id        TUserGroupID `gorm:"primary_key;not null;auto_increment:true"`
	UserId    TUserID      `gorm:"not null;unique_index:group_unique;"`
	GroupId   TGroupID     `gorm:"not null;unique_index:group_unique;"`
	GroupName string       `gorm:"not null"`
	Times

	User *User `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}
