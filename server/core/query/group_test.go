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
			_, err = AddUserGroup(
				db,
				user1.UserId,
				"WICS",
				"Women in Computer Science",
			)
			assert.NoError(t, err)
			_, err = AddUserGroup(
				db,
				user3.UserId,
				"ENG_MENTORSHIP",
				"Engineering Mentorship",
			)
			assert.NoError(t, err)

			err = CreateUserGroups(
				db,
				nil,
				[]data.TUserID{user1.UserId, user2.UserId},
				"WICS",
				"Women in Computer Science",
			)

			users, err := GetUsersByGroupId(db, "WICS")
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
			_, err = AddUserGroup(
				db,
				user1.UserId,
				"WICS",
				"Women in Computer Science",
			)
			assert.NoError(t, err)
			_, err = AddUserGroup(
				db,
				user2.UserId,
				"ENG_MENTORSHIP",
				"Engineering Mentorship",
			)
			assert.NoError(t, err)

			userGroups, err := GetUserGroups(db, user1.UserId)
			assert.NoError(t, err)
			assert.Equal(t, len(userGroups), 1)

			assert.Equal(t, userGroups[0].UserId, user1.UserId)
			assert.Equal(t, userGroups[0].GroupId, data.TGroupID("WICS"))
			assert.Equal(t, userGroups[0].GroupName, "Women in Computer Science")
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

			err = CreateUserGroups(
				db,
				nil,
				[]data.TUserID{user1.UserId, user1.UserId + 1},
				"WICS",
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

			userGroup, err := AddUserGroup(db, user1.UserId, "WICS", "Women in Computer Science")
			assert.NoError(t, err)
			assert.Equal(t, userGroup.UserId, user1.UserId)
			assert.Equal(t, userGroup.GroupId, data.TGroupID("WICS"))
			assert.Equal(t, userGroup.GroupName, "Women in Computer Science")

			_, err = AddUserGroup(db, user1.UserId, "WICS", "Women in Computer Science")
			assert.Error(t, err)
			assert.Equal(
				t,
				err.Error(),
				fmt.Sprintf("You are already a part of the Women in Computer Science group"),
			)
		},
		TestName: "Tests add user group and that we return an error if user group already exists",
	}
	test.RunTestWithDb(thisTest)
}

func TestRemoveUserGroup(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)

			userGroup, err := AddUserGroup(db, user1.UserId, "WICS", "Women in Computer Science")
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
