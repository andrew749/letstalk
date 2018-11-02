package query

import (
	"fmt"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
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
	es *elastic.Client,
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
	var userGroupToIndex *data.UserGroup = nil
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
			userGroupToIndex = userGroup
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}
	if userGroupToIndex != nil {
		// Doesn't crash the program but will print errors to stdout. If this fails, you should
		// run the multi_trait_backfill_es
		indexGroupMultiTrait(es, userGroupToIndex)
	}
	return nil
}

func GetUserGroups(db *gorm.DB, userId data.TUserID) ([]data.UserGroup, errs.Error) {
	var userGroups []data.UserGroup
	err := db.Where(&data.UserGroup{UserId: userId}).Find(&userGroups).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	return userGroups, nil
}
