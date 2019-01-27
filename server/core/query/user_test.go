package query

import (
	"testing"
	"time"

	"letstalk/server/core/test"
	"letstalk/server/core/utility/uw_email"
	"letstalk/server/test_helpers"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSearchUserFallback(t *testing.T) {
	test.RunTestsWithDb([]test.Test{
		{
			Test: func(db *gorm.DB) {
				user, err := test_helpers.CreateTestSetupUser(db, 1)
				assert.NoError(t, err)

				watEmail := uw_email.FromString("test@uwaterloo.ca")
				_, err = GenerateNewVerifyEmailId(db, user.UserId, watEmail)
				assert.NoError(t, err)

				user2, err := GetUserByEmail(db, watEmail.ToStringRaw())
				assert.NoError(t, err)
				assert.Equal(t, user.FirstName, user2.FirstName)
				assert.Equal(t, user.Email, user2.Email)
			},
		},
		{
			TestName: "Test return nil on bad email",
			Test: func(db *gorm.DB) {
				_, err := test_helpers.CreateTestSetupUser(db, 2)
				assert.NoError(t, err)

				badEmail := "invalid_email@uwaterloo.ca"
				user2, err := GetUserByEmail(db, badEmail)
				assert.NoError(t, err)
				assert.Nil(t, user2)
			},
		},
	})
}

func TestGetUsersByCreatedAt(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)
			user2, err := test_helpers.CreateTestSetupUser(db, 2)
			assert.NoError(t, err)
			user3, err := test_helpers.CreateTestSetupUser(db, 3)
			assert.NoError(t, err)

			now := time.Now()

			user1.CreatedAt = now.AddDate(0, 0, 2)
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user3.CreatedAt = now.AddDate(0, 0, -2)
			err = db.Save(user3).Error
			assert.NoError(t, err)

			from := now.AddDate(0, 0, -1)
			to := now.AddDate(0, 0, 1)
			users, err := GetUsersByCreatedAt(db, &from, &to)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, user2.UserId, users[0].UserId)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUsersByCreatedAtBoundaryInclusive(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)
			_, err = test_helpers.CreateTestSetupUser(db, 2)
			assert.NoError(t, err)
			user3, err := test_helpers.CreateTestSetupUser(db, 3)
			assert.NoError(t, err)

			now := time.Now()

			user1.CreatedAt = now.AddDate(0, 0, 1)
			err = db.Save(user1).Error
			assert.NoError(t, err)

			user3.CreatedAt = now.AddDate(0, 0, -1)
			err = db.Save(user3).Error
			assert.NoError(t, err)

			from := now.AddDate(0, 0, -1)
			to := now.AddDate(0, 0, 1)
			users, err := GetUsersByCreatedAt(db, &from, &to)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(users))
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUsersByCreatedAtNotSpecified(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			for i := 0; i < 3; i++ {
				_, err := test_helpers.CreateTestSetupUser(db, i+1)
				assert.NoError(t, err)
			}

			users, err := GetUsersByCreatedAt(db, nil, nil)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(users))
		},
	}
	test.RunTestWithDb(thisTest)
}
