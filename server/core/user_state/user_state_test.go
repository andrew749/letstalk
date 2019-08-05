package user_state

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/api"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
)

func TestGetUserStateNoSuchUser(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			state, err := GetUserState(db, 1)
			assert.Error(t, err)
			assert.Nil(t, state)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUserStateAccountCreated(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)

			state, err := GetUserState(db, user.UserId)
			assert.NoError(t, err)
			assert.Equal(t, api.ACCOUNT_CREATED, *state)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUserStateAccountEmailVerified(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)

			user.IsEmailVerified = true
			err = db.Save(user).Error
			assert.NoError(t, err)

			state, err := GetUserState(db, user.UserId)
			assert.NoError(t, err)
			assert.Equal(t, api.ACCOUNT_EMAIL_VERIFIED, *state)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUserStateAccountHasBasicInfo(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)

			user.IsEmailVerified = true
			err = db.Save(user).Error
			assert.NoError(t, err)

			cohort := &data.Cohort{
				ProgramId:   "ARTS",
				ProgramName: "Arts",
				GradYear:    2018,
				IsCoop:      false,
			}
			err = db.Save(cohort).Error
			assert.NoError(t, err)

			userCohort := &data.UserCohort{
				UserId:   user.UserId,
				CohortId: cohort.CohortId,
			}
			err = db.Save(userCohort).Error
			assert.NoError(t, err)

			bio := "Hello"
			userAdditionalData := &data.UserAdditionalData{
				Bio: &bio,
			}
			err = db.Save(userAdditionalData).Error
			assert.NoError(t, err)

			state, err := GetUserState(db, user.UserId)
			assert.NoError(t, err)
			assert.Equal(t, api.ACCOUNT_HAS_BASIC_INFO, *state)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUserStateAccountSetup(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user, err := test_helpers.CreateTestSetupUser(db, 1)
			assert.NoError(t, err)

			state, err := GetUserState(db, user.UserId)
			assert.NoError(t, err)
			assert.Equal(t, api.ACCOUNT_SETUP, *state)
		},
	}
	test.RunTestWithDb(thisTest)
}
