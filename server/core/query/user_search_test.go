package query

import (
	"fmt"
	"testing"

	"letstalk/server/core/api"
	"letstalk/server/core/test"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func createUser(db *gorm.DB, num int) (*data.User, error) {
	cohort := &data.Cohort{
		ProgramId:   "ARTS",
		ProgramName: "Arts",
		GradYear:    2018 + uint(num),
		IsCoop:      false,
	}
	err := db.Save(cohort).Error
	if err != nil {
		return nil, err
	}

	user, err := data.CreateUser(
		db,
		fmt.Sprintf("john.doe%d@gmail.com", num),
		fmt.Sprintf("John%d", num),
		fmt.Sprintf("Doe%d", num),
		data.GENDER_MALE,
		"1996-11-07",
		data.USER_ROLE_DEFAULT,
	)
	if err != nil {
		return nil, err
	}

	userCohort := &data.UserCohort{
		UserId:   user.UserId,
		CohortId: cohort.CohortId,
	}
	err = db.Save(userCohort).Error
	if err != nil {
		return nil, err
	}
	userCohort.Cohort = cohort
	user.Cohort = userCohort
	return user, nil
}

func assertUserSearchResultEqual(t *testing.T, res api.UserSearchResult, user *data.User) {
	assert.Equal(t, user.UserId, res.UserId)
	assert.Equal(t, user.FirstName, res.FirstName)
	assert.Equal(t, user.LastName, res.LastName)
	assert.Equal(t, user.Gender, res.Gender)
	assert.NotNil(t, res.Cohort)
	assert.Equal(t, user.Cohort.Cohort.CohortId, res.Cohort.CohortId)
	assert.Equal(t, user.Cohort.Cohort.ProgramId, res.Cohort.ProgramId)
	assert.Equal(t, user.Cohort.Cohort.ProgramName, res.Cohort.ProgramName)
	assert.Equal(t, user.Cohort.Cohort.GradYear, res.Cohort.GradYear)
	assert.Equal(t, user.Cohort.Cohort.IsCoop, res.Cohort.IsCoop)
	assert.Equal(t, user.Cohort.Cohort.SequenceId, res.Cohort.SequenceId)
	assert.Equal(t, user.Cohort.Cohort.SequenceName, res.Cohort.SequenceName)
	assert.Equal(t, user.ProfilePic, res.ProfilePic)
	assert.Nil(t, res.Reason)
}

func TestSearchUsersByCohort(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error

			user1, err := createUser(db, 1)
			assert.NoError(t, err)

			_, err = createUser(db, 2)
			assert.NoError(t, err)

			myUser, err := createUser(db, 3)
			assert.NoError(t, err)

			myUser.Cohort = user1.Cohort
			err = db.Save(myUser.Cohort).Error
			assert.NoError(t, err)

			req := api.CohortUserSearchRequest{
				CohortId:                user1.Cohort.CohortId,
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 10},
			}
			res, err := SearchUsersByCohort(db, req, myUser.UserId)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 1, res.NumResults)
			assert.Equal(t, 1, len(res.Results))

			assertUserSearchResultEqual(t, res.Results[0], user1)
		},
		TestName: "Test correct results when searching for users by cohort",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersByCohortLimit(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			numUsers := 10
			users := make([]data.User, numUsers)
			for i := 0; i < numUsers; i = i + 1 {
				user, err := createUser(db, i+1)
				assert.NoError(t, err)
				users[i] = *user
				userTrait := data.UserCohort{
					UserId:   user.UserId,
					CohortId: users[0].Cohort.CohortId,
				}
				err = db.Save(&userTrait).Error
				assert.NoError(t, err)
			}

			req := api.CohortUserSearchRequest{
				CohortId:                users[0].Cohort.CohortId,
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 5},
			}
			res, err := SearchUsersByCohort(db, req, 69)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 5, res.NumResults)
			assert.Equal(t, 5, len(res.Results))
		},
		TestName: "Test truncating of results when searching for users by cohort",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersByPosition(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error

			user1, err := createUser(db, 1)
			assert.NoError(t, err)

			user2, err := createUser(db, 2)
			assert.NoError(t, err)

			myUser, err := createUser(db, 3)
			assert.NoError(t, err)

			userPosition1 := data.UserPosition{
				UserId:         user1.UserId,
				RoleId:         data.TRoleID(69),
				OrganizationId: data.TOrganizationID(69),
			}
			err = db.Save(&userPosition1).Error
			assert.NoError(t, err)

			// Testing reduping
			userPosition1.Id = 0
			err = db.Save(&userPosition1).Error
			assert.NoError(t, err)

			userPosition2 := data.UserPosition{
				UserId:         user2.UserId,
				RoleId:         data.TRoleID(69),
				OrganizationId: data.TOrganizationID(70),
			}
			err = db.Save(&userPosition2).Error
			assert.NoError(t, err)

			myUserPosition := data.UserPosition{
				UserId:         myUser.UserId,
				RoleId:         data.TRoleID(69),
				OrganizationId: data.TOrganizationID(69),
			}
			err = db.Save(&myUserPosition).Error
			assert.NoError(t, err)

			req := api.PositionUserSearchRequest{
				RoleId:                  data.TRoleID(69),
				OrganizationId:          data.TOrganizationID(69),
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 10},
			}
			res, err := SearchUsersByPosition(db, req, myUser.UserId)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 1, res.NumResults)
			assert.Equal(t, 1, len(res.Results))

			assertUserSearchResultEqual(t, res.Results[0], user1)
		},
		TestName: "Test correct results when searching for users by position",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersByPositionLimit(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			numUsers := 10
			users := make([]data.User, numUsers)
			for i := 0; i < numUsers; i = i + 1 {
				user, err := createUser(db, i+1)
				assert.NoError(t, err)
				users[i] = *user
				userPosition := data.UserPosition{
					UserId:         user.UserId,
					RoleId:         data.TRoleID(69),
					OrganizationId: data.TOrganizationID(69),
				}
				err = db.Save(&userPosition).Error
				assert.NoError(t, err)
			}

			req := api.PositionUserSearchRequest{
				RoleId:                  data.TRoleID(69),
				OrganizationId:          data.TOrganizationID(69),
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 5},
			}
			res, err := SearchUsersByPosition(db, req, 69)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 5, res.NumResults)
			assert.Equal(t, 5, len(res.Results))
		},
		TestName: "Test truncating of results when searching for users by position",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersBySimpleTrait(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error

			user1, err := createUser(db, 1)
			assert.NoError(t, err)

			user2, err := createUser(db, 2)
			assert.NoError(t, err)

			myUser, err := createUser(db, 3)
			assert.NoError(t, err)

			userTrait1 := data.UserSimpleTrait{
				UserId:                 user1.UserId,
				SimpleTraitId:          data.TSimpleTraitID(69),
				SimpleTraitIsSensitive: false,
			}
			err = db.Save(&userTrait1).Error
			assert.NoError(t, err)

			userTrait1.Id = 0
			err = db.Save(&userTrait1).Error
			assert.NoError(t, err)

			userTrait2 := data.UserSimpleTrait{
				UserId:        user2.UserId,
				SimpleTraitId: data.TSimpleTraitID(70),
			}
			err = db.Save(&userTrait2).Error
			assert.NoError(t, err)

			myUserTrait := data.UserSimpleTrait{
				UserId:        myUser.UserId,
				SimpleTraitId: data.TSimpleTraitID(69),
			}
			err = db.Save(&myUserTrait).Error
			assert.NoError(t, err)

			req := api.SimpleTraitUserSearchRequest{
				SimpleTraitId:           data.TSimpleTraitID(69),
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 10},
			}
			res, err := SearchUsersBySimpleTrait(db, req, myUser.UserId)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 1, res.NumResults)
			assert.Equal(t, 1, len(res.Results))

			assertUserSearchResultEqual(t, res.Results[0], user1)
		},
		TestName: "Test correct results when searching for users by simple trait",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersBySimpleTraitAnon(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error

			user1, err := createUser(db, 1)
			assert.NoError(t, err)

			user2, err := createUser(db, 2)
			assert.NoError(t, err)

			userTrait1 := data.UserSimpleTrait{
				UserId:                 user1.UserId,
				SimpleTraitId:          data.TSimpleTraitID(69),
				SimpleTraitIsSensitive: true,
			}
			err = db.Save(&userTrait1).Error
			assert.NoError(t, err)

			userTrait2 := data.UserSimpleTrait{
				UserId:        user2.UserId,
				SimpleTraitId: data.TSimpleTraitID(70),
			}
			err = db.Save(&userTrait2).Error
			assert.NoError(t, err)

			req := api.SimpleTraitUserSearchRequest{
				SimpleTraitId:           data.TSimpleTraitID(69),
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 10},
			}
			res, err := SearchUsersBySimpleTrait(db, req, 69)
			assert.NoError(t, err)
			assert.Equal(t, true, res.IsAnonymous)
			assert.Equal(t, 1, res.NumResults)
			assert.Equal(t, 0, len(res.Results))
		},
		TestName: "Test correct results when searching for users by simple trait with anon",
	}
	test.RunTestWithDb(thisTest)
}

func TestSearchUsersBySimpleTraitLimit(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			numUsers := 10
			users := make([]data.User, numUsers)
			for i := 0; i < numUsers; i = i + 1 {
				user, err := createUser(db, i+1)
				assert.NoError(t, err)
				users[i] = *user
				userTrait := data.UserSimpleTrait{
					UserId:                 user.UserId,
					SimpleTraitId:          data.TSimpleTraitID(69),
					SimpleTraitIsSensitive: false,
				}
				err = db.Save(&userTrait).Error
				assert.NoError(t, err)
			}

			req := api.SimpleTraitUserSearchRequest{
				SimpleTraitId:           data.TSimpleTraitID(69),
				CommonUserSearchRequest: api.CommonUserSearchRequest{Size: 5},
			}
			res, err := SearchUsersBySimpleTrait(db, req, 69)
			assert.NoError(t, err)
			assert.Equal(t, false, res.IsAnonymous)
			assert.Equal(t, 5, res.NumResults)
			assert.Equal(t, 5, len(res.Results))
		},
		TestName: "Test truncating of results when searching for users by simple trait",
	}
	test.RunTestWithDb(thisTest)
}
