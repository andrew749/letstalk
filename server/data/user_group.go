package data

import "github.com/jinzhu/gorm"

type (
	TUserGroupID  EntID
	TGroupID      string
	UserGroupType string
)

// Stores groups that a user is a part of
type UserGroup struct {
	Id        TUserGroupID `gorm:"primary_key;not null;auto_increment:true"`
	UserId    TUserID      `gorm:"not null"`
	GroupId   TGroupID     `gorm:"not null"`
	GroupName string       `gorm:"not null"`
	Times

	User *User `gorm:"foreignkey:UserId;association_foreignkey:UserId"`
}

func CreateUserGroup(
	db *gorm.DB,
	userId TUserID,
	groupId TGroupID,
	groupName string,
) (*UserGroup, error) {
	userGroup := &UserGroup{
		UserId:    userId,
		GroupId:   groupId,
		GroupName: groupName,
	}
	err := db.Create(userGroup).Error
	if err != nil {
		return nil, err
	}
	return userGroup, nil
}
