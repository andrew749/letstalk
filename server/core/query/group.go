package query

import (
	"fmt"

	"github.com/google/uuid"

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
			userGroup := data.UserGroup{
				UserId:    userId,
				GroupId:   groupId,
				GroupName: groupName,
			}
			err := db.Where(&userGroup).FirstOrCreate(&userGroup).Error
			if err != nil {
				rlog.Errorf("Failed on user %d: %v\n", userId, err)
				return err
			}
			userGroupToIndex = &userGroup
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

// CreateGroup Create a new group.
func CreateGroup(db *gorm.DB, groupName string) (*data.Group, errs.Error) {
	groupUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}
	group := data.Group{GroupId: data.TGroupID(groupUUID.String()), GroupName: groupName}
	err = db.Create(&group).Error
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}
	return &group, nil
}

// GetUserGroups Find all groups a user is part of.
func GetUserGroups(db *gorm.DB, userId data.TUserID) ([]data.UserGroup, errs.Error) {
	var userGroups []data.UserGroup
	err := db.Where(&data.UserGroup{UserId: userId}).Find(&userGroups).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	return userGroups, nil
}

// AddUserGroup Add a user to a group. Idempotent operation.
func AddUserGroup(
	db *gorm.DB,
	userId data.TUserID,
	groupId data.TGroupID,
) (*data.UserGroup, errs.Error) {
	userGroup := data.UserGroup{
		UserId:  userId,
		GroupId: groupId,
	}
	var group data.Group
	res := db.Where(&data.Group{GroupId: groupId}).First(&group)
	if res.RecordNotFound() {
		return nil, errs.NewRequestError("Invalid group id.")
	}

	if err := db.Where(&data.UserGroup{UserId: userId, GroupId: groupId, GroupName: group.GroupName}).FirstOrCreate(&userGroup).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	return &userGroup, nil
}

func RemoveUserGroup(db *gorm.DB, userId data.TUserID, userGroupId data.TUserGroupID) errs.Error {
	userGroup := data.UserGroup{
		Id:     userGroupId,
		UserId: userId,
	}
	if err := db.Where(&userGroup).Delete(&data.UserGroup{}).Error; err != nil {
		return errs.NewDbError(err)
	}
	return nil
}

// RemoveUserFromGroup Removes a user from a given group.
func RemoveUserFromGroup(db *gorm.DB, userId data.TUserID, groupId data.TGroupID) errs.Error {
	userGroup, err := GetUserGroupForUserIdGroupId(db, userId, groupId)
	if err != nil {
		return err
	}

	return RemoveUserGroup(db, userId, userGroup.Id)
}

// GetUserGroupForUserIdGroupId get a UserGroup object given user id and group id
func GetUserGroupForUserIdGroupId(db *gorm.DB, userId data.TUserID, groupId data.TGroupID) (*data.UserGroup, errs.Error) {
	var res data.UserGroup
	if err := db.Where(&data.UserGroup{UserId: userId, GroupId: groupId}).First(&res).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	return &res, nil
}

// EnrollUserInManagedGroup Enroll the user into an administrator managed group.
func EnrollUserInManagedGroup(db *gorm.DB, userId data.TUserID, groupId data.TGroupID) errs.Error {
	var managedGroup data.ManagedGroup

	if groupId == "" {
		return errs.NewRequestError("Group not provided")
	}

	// find the group by uuid
	if err := db.Where(&data.ManagedGroup{GroupId: groupId}).Preload("Group").First(&managedGroup).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errs.NewRequestError("Group does not exist")
		}

		return errs.NewInternalError(err.Error())
	}
	// The group exists for the mapping
	_, err := AddUserGroup(db, userId, managedGroup.Group.GroupId)
	return err
}

// CreateManagedGroup Create a group that this admin manages.
func CreateManagedGroup(
	db *gorm.DB,
	adminUserID data.TUserID,
	groupName string,
) (*data.ManagedGroup, errs.Error) {
	var managedGroup data.ManagedGroup
	group, err := CreateGroup(db, groupName)

	if err != nil {
		return nil, err
	}

	managedGroup.GroupId = group.GroupId
	managedGroup.AdministratorId = adminUserID

	err2 := db.Create(&managedGroup).Error
	if err2 != nil {
		return nil, errs.NewBaseError(err2.Error())
	}

	return &managedGroup, nil
}

func CheckAdminManagesGroup(db *gorm.DB, userId data.TUserID, groupId data.TGroupID) (bool, errs.Error) {
	var res data.ManagedGroup
	if err := db.Where(&data.ManagedGroup{
		AdministratorId: userId,
		GroupId:         groupId,
	}).First(&res).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		} else {
			return false, errs.NewDbError(err)
		}
	}
	return true, nil
}

// GetManagedGroups Get all the groups that the admin manages.
func GetManagedGroups(
	db *gorm.DB,
	adminUserID data.TUserID,
) ([]data.ManagedGroup, errs.Error) {
	var groups []data.ManagedGroup
	if err := db.Where(&data.ManagedGroup{AdministratorId: adminUserID}).Preload("Group").Find(&groups).Error; err != nil {
		return nil, errs.NewInternalError(err.Error())
	}
	rlog.Infof("%+v", groups)
	return groups, nil
}
