package query

import (
	"fmt"
	"testing"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByGroupId(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)
			user2, err := test_helpers.CreateTestSetupUser(db, 2)
			assert.NoError(t, err)
			user3, err := test_helpers.CreateTestSetupUser(db, 3)
			assert.NoError(t, err)
			wicsGroup, err := CreateGroup(db, "WICS")
			assert.NoError(t, err)
			engGroup, err := CreateGroup(db, "ENG_MENTORSHIP")
			assert.NoError(t, err)
			_, err = AddUserGroup(
				db,
				user3.UserId,
				engGroup.GroupId,
			)
			assert.NoError(t, err)

			err = CreateUserGroups(
				db,
				nil,
				[]data.TUserID{user1.UserId, user2.UserId},
				wicsGroup.GroupId,
				wicsGroup.GroupName,
			)

			users, err := GetUsersByGroupId(db, wicsGroup.GroupId)
			assert.NoError(t, err)
			assert.Equal(t, len(users), 2)

			assert.Equal(t, users[0].UserId, user1.UserId)
			assert.Equal(t, users[1].UserId, user2.UserId)
		},
		TestName: "Test get users by group id",
	}
	test.RunTestWithDb(thisTest)
}

func TestUserGroups(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)
			user2, err := test_helpers.CreateTestSetupUser(db, 2)
			assert.NoError(t, err)
			wicsGroup, err := CreateGroup(db, "WICS")
			assert.NoError(t, err)
			engGroup, err := CreateGroup(db, "ENG_MENTORSHIP")
			assert.NoError(t, err)
			_, err = AddUserGroup(
				db,
				user1.UserId,
				wicsGroup.GroupId,
			)
			assert.NoError(t, err)
			_, err = AddUserGroup(
				db,
				user2.UserId,
				engGroup.GroupId,
			)
			assert.NoError(t, err)

			userGroups, err := GetUserGroups(db, user1.UserId)
			assert.NoError(t, err)
			assert.Equal(t, len(userGroups), 1)

			assert.Equal(t, userGroups[0].UserId, user1.UserId)
			assert.Equal(t, userGroups[0].GroupId, wicsGroup.GroupId)
			assert.Equal(t, userGroups[0].GroupName, "WICS")
		},
		TestName: "Test get user groups by user id",
	}
	test.RunTestWithDb(thisTest)
}

func TestCreateUserGroupsMissingUsers(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)

			group, err := CreateGroup(db, "WICS")
			assert.NoError(t, err)
			err = CreateUserGroups(
				db,
				nil,
				[]data.TUserID{user1.UserId, user1.UserId + 1},
				group.GroupId,
				"Women in Computer Science",
			)
			assert.Error(t, err)
			assert.Equal(
				t,
				err.Error(),
				fmt.Sprintf("Missing users: %v", []data.TUserID{user1.UserId + 1}),
			)
		},
		TestName: "Test get users by group id fails when some users don't exist",
	}
	test.RunTestWithDb(thisTest)
}

func TestAddUserGroup(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)

			group, err := CreateGroup(db, "WICS")
			assert.NoError(t, err)

			userGroup, err := AddUserGroup(db, user1.UserId, group.GroupId)
			assert.NoError(t, err)
			assert.Equal(t, userGroup.UserId, user1.UserId)
			assert.Equal(t, userGroup.GroupName, "WICS")

			_, err = AddUserGroup(db, user1.UserId, group.GroupId)
			assert.NoError(t, err)
		},
		TestName: "Tests add user to group and that we don't return an error if user is already part of the group",
	}
	test.RunTestWithDb(thisTest)
}

func TestRemoveUserGroup(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)

			group, err := CreateGroup(db, "WICS")
			assert.NoError(t, err)

			userGroup, err := AddUserGroup(db, user1.UserId, group.GroupId)
			assert.NoError(t, err)

			var userGroups []data.UserGroup
			err = db.Where(&data.UserGroup{UserId: user1.UserId}).Find(&userGroups).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(userGroups))

			err = RemoveUserGroup(db, user1.UserId, userGroup.Id)
			assert.NoError(t, err)

			err = db.Where(&data.UserGroup{UserId: user1.UserId}).Find(&userGroups).Error
			assert.NoError(t, err)
			assert.Equal(t, 0, len(userGroups))
		},
		TestName: "Test that removing user groups works",
	}
	test.RunTestWithDb(thisTest)
}
