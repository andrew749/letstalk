package query

import (
	"fmt"
	"testing"

	"letstalk/server/core/test"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByGroupId(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := createTestUser(db, 1)
			assert.NoError(t, err)
			user2, err := createTestUser(db, 2)
			assert.NoError(t, err)
			user3, err := createTestUser(db, 3)
			assert.NoError(t, err)
			_, err = data.CreateUserGroup(db, user1.UserId, "WICS", "Women in Computer Science")
			assert.NoError(t, err)
			_, err = data.CreateUserGroup(db, user3.UserId, "ENG_MENTORSHIP", "Engineering Mentorship")
			assert.NoError(t, err)

			err = CreateUserGroups(
				db,
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

func TestCreateUserGroupsMissingUsers(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := createTestUser(db, 1)
			assert.NoError(t, err)

			err = CreateUserGroups(
				db,
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
