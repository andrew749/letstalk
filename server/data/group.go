package data

import "github.com/jinzhu/gorm"

// ManagedGroup A group that is managed by an administrator.
// application logic ensures that group name is unique since we use this as a unique key in other application logic.
type ManagedGroup struct {
	gorm.Model
	AdministratorId TUserID
	Group           Group `gorm:"foreignkey:GroupId;association_foreignkey:GroupId;"`
	GroupId         TGroupID
}

type Group struct {
	Times
	GroupId   TGroupID `gorm:"primary_key;"`
	GroupName string   `gorm:"not null;size:100"`
}
