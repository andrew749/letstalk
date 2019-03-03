package generic_notification_job

import (
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/test_helpers"
	"testing"
	"time"

	"github.com/romana/rlog"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var specStore = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		GENERIC_NOTIFICATION_JOB: GenericNotificationJob,
	},
}

func createUsers(db *gorm.DB) (map[int]data.User, error) {
	sequenceId := "8_STREAM"
	sequenceName := "8 Stream"
	now := time.Now()

	user1, err := test_helpers.CreateTestUser(db, 1)
	if err != nil {
		return nil, err
	}
	user1.CreatedAt = now.AddDate(0, 0, -2)
	user1.Gender = data.GENDER_FEMALE
	err = db.Save(user1).Error
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user2, err := test_helpers.CreateTestUser(db, 2)
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user3, err := test_helpers.CreateTestUser(db, 3)
	if err != nil {
		return nil, err
	}
	user3.CreatedAt = now.AddDate(0, 0, 2)
	err = db.Save(user3).Error
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user3, "COMPUTER_ENGINEERING", "Computer Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	return map[int]data.User{1: *user1, 2: *user2, 3: *user3}, nil
}

func nukeUsers(db *gorm.DB) error {
	return db.Exec("DELETE FROM users;").Error
}

// Because we bind parameters at the end we can't really test the notifications or emails produced.
// for now we just test user selection.
func TestGetMetaDataForQuerySimple(t *testing.T) {
	theseTests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				users, err := createUsers(db)
				assert.NoError(t, err)
				metadata, err := getMetadataForQuery(db, "SELECT * from users where user_id=1;")
				assert.NoError(t, err)
				assert.Len(t, metadata, 1)
				assert.Equal(t, data.TUserID(1), metadata[0].UserId)
				assert.Equal(t, users[1].FirstName, metadata[0].Data["first_name"])
			},
		},
	}
	test.RunTestsWithDb(theseTests)

}

func TestGetMetaDataForQueryJoin(t *testing.T) {
	theseTests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				users, err := createUsers(db)
				assert.NoError(t, err)
				metadata, err := getMetadataForQuery(db, "SELECT u.user_id as user_id, cohort_id from users u, user_cohorts uc where u.user_id=uc.user_id order by u.user_id;")
				assert.NoError(t, err)
				assert.Len(t, metadata, 3)
				assert.Equal(t, users[1].UserId, metadata[0].UserId)
				assert.Equal(t, users[1].Cohort.CohortId, data.TCohortID(metadata[0].Data["cohort_id"].(int64)))
				assert.Equal(t, users[2].UserId, metadata[1].UserId)
				assert.Equal(t, users[2].Cohort.CohortId, data.TCohortID(metadata[1].Data["cohort_id"].(int64)))
				assert.Equal(t, users[3].UserId, metadata[2].UserId)
				assert.Equal(t, users[3].Cohort.CohortId, data.TCohortID(metadata[2].Data["cohort_id"].(int64)))
				nukeUsers(db)
			},
		},
	}
	test.RunTestsWithDb(theseTests)

}

func TestGetMetaDataForQueryError(t *testing.T) {
	theseTests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				// errors
				_, err := getMetadataForQuery(db, "DROP")
				rlog.Debugf("%+v", err)
				assert.Error(t, err)
				_, err = getMetadataForQuery(db, "DELETE")
				assert.Error(t, err)
				_, err = getMetadataForQuery(db, "INSERT")
				assert.Error(t, err)
				_, err = getMetadataForQuery(db, "CREATE")
				assert.Error(t, err)
				_, err = getMetadataForQuery(db, "update")
				assert.Error(t, err)
			},
		},
	}
	test.RunTestsWithDb(theseTests)

}

func TestMergeMetadata(t *testing.T) {
	map1 := map[string]interface{}{"a": "b", "b": 1}
	map2 := map[string]interface{}{"c": "d", "e": 1}
	map3 := map[string]interface{}{"a": "f"}

	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				// test merging normal maps
				mergedMaps, err := mergeMaps(map1, map2)
				assert.NoError(t, err)
				assert.Exactly(t, map[string]interface{}{"a": "b", "b": 1, "c": "d", "e": 1}, mergedMaps)
			},
		},
		test.Test{
			Test: func(db *gorm.DB) {
				// test merging maps with duplicate key
				_, err := mergeMaps(map1, map3)
				assert.Error(t, err)
			},
		},
	}
	test.RunTestsWithDb(tests)
}
