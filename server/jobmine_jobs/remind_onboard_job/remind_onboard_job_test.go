package remind_onboard_job

import (
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// test that the methods get users
func TestDataSearching(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				var err error
				user1Id := 1
				user2Id := 2
				user, err := query.CreateTestUser(db, user1Id)
				assert.NoError(t, err)
				user2, err := query.CreateTestUser(db, user2Id)
				assert.NoError(t, err)
				_, err = query.AddUserSimpleTraitByName(db, nil, data.TUserID(user1Id), "test trait")
				assert.NoError(t, err)
				users, err := usersWithoutTraits(db)
				assert.NoError(t, err)
				assert.NotContains(t, *users, user.UserId)
				assert.Contains(t, *users, user2.UserId)
			},
			TestName: "Test find users who don't have any traits",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				var err error
				user1Id := 3
				user2Id := 4
				user, err := query.CreateTestUser(db, user1Id)
				assert.NoError(t, err)

				user2, err := query.CreateTestUser(db, user2Id)
				assert.NoError(t, err)

				bio := "Give me the zucc"
				err = db.Save(&data.UserAdditionalData{UserId: data.TUserID(user1Id), Bio: &bio}).Error
				assert.NoError(t, err)

				users, err := usersWithoutBio(db)
				assert.NoError(t, err)

				assert.NotContains(t, *users, user.UserId)
				assert.Contains(t, *users, user2.UserId)
			},
			TestName: "Test find users who don't have a bio",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				var err error
				user1Id := 5
				user2Id := 6
				user, err := query.CreateTestUser(db, user1Id)
				assert.NoError(t, err)
				user2, err := query.CreateTestUser(db, user2Id)
				assert.NoError(t, err)

				positionName := "test"
				organizationName := "test"
				_, err = query.AddUserPosition(db, nil, data.TUserID(user1Id), nil, &positionName, nil, &organizationName, "1996-10-07", nil)
				assert.NoError(t, err)
				users, err := usersWithoutPosition(db)
				assert.NoError(t, err)

				assert.NotContains(t, *users, user.UserId)
				assert.Contains(t, *users, user2.UserId)
			},
			TestName: "Test find users who don't have a position",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				var err error
				user1Id := 7
				user2Id := 8
				user, err := query.CreateTestUser(db, user1Id)
				assert.NoError(t, err)
				user2, err := query.CreateTestUser(db, user2Id)
				assert.NoError(t, err)

				err = query.CreateUserGroups(db, nil, []data.TUserID{data.TUserID(user1Id)}, data.TGroupID("test"), "test")
				assert.NoError(t, err)

				users, err := usersWithoutGroup(db)
				assert.NoError(t, err)

				assert.NotContains(t, *users, user.UserId)
				assert.Contains(t, *users, user2.UserId)
			},
			TestName: "Test find users who don't have a group",
		},
	}
	test.RunTestsWithDb(tests)
}
