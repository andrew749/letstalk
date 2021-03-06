package seed_mentorships_job

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
)

// TODO: Can consolidate these tests
func TestGetLowerYears(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "MECHATRONICS_ENGINEERING", "Mechatronics Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user4, err := test_helpers.CreateTestUser(db, 4)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user4, "ENVIRONMENTAL_ENGINEERING", "Environmental Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			userIds, err := GetCandidates(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022},
				true, 2021, nil, nil)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user1.UserId, user3.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetUpperYears(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "MECHATRONICS_ENGINEERING", "Mechatronics Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user4, err := test_helpers.CreateTestUser(db, 4)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user4, "ENVIRONMENTAL_ENGINEERING", "Environmental Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			userIds, err := GetCandidates(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022}, false,
				2021, nil, nil)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user2.UserId, user3.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetLowerYearsWithinCreatedAtRange(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"
			now := time.Now()

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.CreatedAt = now.AddDate(0, 0, -2)
			err = db.Save(user1).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			user3.CreatedAt = now.AddDate(0, 0, 2)
			err = db.Save(user3).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			from := now.AddDate(0, 0, -1)
			to := now.AddDate(0, 0, 1)
			userIds, err := GetCandidates(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022}, true,
				2021, &from, &to)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user2.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetLowerYearsWithinCreatedAtRangeBoundaryInclusive(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"
			now := time.Now()

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.CreatedAt = now.AddDate(0, 0, -2)
			err = db.Save(user1).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			user3.CreatedAt = now.AddDate(0, 0, 2)
			err = db.Save(user3).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			from := now.AddDate(0, 0, -2)
			to := now.AddDate(0, 0, 2)
			userIds, err := GetCandidates(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022}, true,
				2021, &from, &to)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user1.UserId, user2.UserId, user3.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetLowerUpperYears(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "MECHATRONICS_ENGINEERING", "Mechatronics Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user4, err := test_helpers.CreateTestUser(db, 4)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user4, "ENVIRONMENTAL_ENGINEERING", "Environmental Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			userIds, err := GetFilteredLowerAndAllUpperYears(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022}, 2021,
				nil, nil)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user1.UserId, user2.UserId, user3.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetLowerUpperYearsRestrictGradYears(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			userIds, err := GetFilteredLowerAndAllUpperYears(
				db, []string{"SOFTWARE_ENGINEERING"}, []uint{2021}, 2000, nil, nil)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{user2.UserId}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetLowerUpperYearsRanges(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			sequenceId := "8_STREAM"
			sequenceName := "8 Stream"
			now := time.Now()

			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			user1.CreatedAt = now.AddDate(0, 0, -2)
			err = db.Save(user1).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user3, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			user3.CreatedAt = now.AddDate(0, 0, 2)
			err = db.Save(user3).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user3, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user4, err := test_helpers.CreateTestUser(db, 4)
			assert.NoError(t, err)
			user4.CreatedAt = now.AddDate(0, 0, -2)
			err = db.Save(user4).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user4, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user5, err := test_helpers.CreateTestUser(db, 5)
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user5, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			user6, err := test_helpers.CreateTestUser(db, 6)
			assert.NoError(t, err)
			user6.CreatedAt = now.AddDate(0, 0, 2)
			err = db.Save(user6).Error
			assert.NoError(t, err)
			err = test_helpers.CreateCohortForUser(
				db, user6, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
				&sequenceId, &sequenceName)
			assert.NoError(t, err)

			from := now.AddDate(0, 0, -1)
			to := now.AddDate(0, 0, 1)
			userIds, err := GetFilteredLowerAndAllUpperYears(
				db, []string{"SOFTWARE_ENGINEERING", "MECHATRONICS_ENGINEERING"}, []uint{2021, 2022}, 2021,
				&from, &to)
			assert.NoError(t, err)

			assert.ElementsMatch(t, []data.TUserID{
				user2.UserId, user4.UserId, user5.UserId, user6.UserId,
			}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}
