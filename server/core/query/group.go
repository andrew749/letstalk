package query

import (
	"fmt"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

func GetUsersByGroupId(db *gorm.DB, groupId data.TGroupID) ([]data.User, errs.Error) {
	var userGroups []data.UserGroup
	err := db.Where(&data.UserGroup{GroupId: groupId}).Preload("User").Find(&userGroups).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	users := make([]data.User, len(userGroups))
	for i, userGroup := range userGroups {
		users[i] = *userGroup.User
	}
	return users, nil
}

func CreateUserGroups(
	db *gorm.DB,
	userIds []data.TUserID,
	groupId data.TGroupID,
	groupName string,
) errs.Error {
	missingUserIds, err := MissingUsers(db, userIds)
	if err != nil {
		return err
	}
	if len(missingUserIds) != 0 {
		return errs.NewRequestError(fmt.Sprintf("Missing users: %v", missingUserIds))
	}
	dbErr := ctx.WithinTx(db, func(db *gorm.DB) error {
		for _, userId := range userIds {
			userGroup := &data.UserGroup{
				UserId:    userId,
				GroupId:   groupId,
				GroupName: groupName,
			}
			err := db.Where(userGroup).FirstOrCreate(userGroup).Error
			if err != nil {
				rlog.Errorf("Failed on user %d: %v\n", userId, err)
				return err
			}
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	return nil
}
