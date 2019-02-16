package seed_mentorships_job

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/test_helpers"
	"letstalk/server/utility"
)

type match struct {
	mentorId data.TUserID
	menteeId data.TUserID
}

var specStore = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		SEED_MENTORSHIPS_JOB: SeedJobSpec,
	},
}

func createUsers(db *gorm.DB) ([]data.User, error) {
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

	user4, err := test_helpers.CreateTestUser(db, 4)
	if err != nil {
		return nil, err
	}
	user4.CreatedAt = now.AddDate(0, 0, -2)
	user4.Gender = data.GENDER_FEMALE
	err = db.Save(user4).Error
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user4, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user5, err := test_helpers.CreateTestUser(db, 5)
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user5, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user6, err := test_helpers.CreateTestUser(db, 6)
	if err != nil {
		return nil, err
	}
	user6.CreatedAt = now.AddDate(0, 0, 2)
	err = db.Save(user6).Error
	if err != nil {
		return nil, err
	}
	err = test_helpers.CreateCohortForUser(
		db, user6, "COMPUTER_ENGINEERING", "Computer Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}
	return []data.User{*user1, *user2, *user3, *user4, *user5, *user6}, nil
}

func runAndTestRunners(t *testing.T, db *gorm.DB, runId string) {
	now := time.Now()
	var queryJob jobmine.JobRecord
	err := db.Where("run_id = ?", runId).Find(&queryJob).Error
	assert.NoError(t, err)
	assert.NotZero(t, queryJob.ID)
	assert.Equal(t, jobmine.STATUS_CREATED, queryJob.Status)

	// set start time a little earlier
	queryJob.StartTime = now.AddDate(0, -1, 0) // yesterday
	err = db.Save(&queryJob).Error
	assert.NoError(t, err)

	// schedule tasks
	err = jobmine.JobRunner(specStore, db)
	assert.NoError(t, err)

	// run tasks
	err = jobmine.TaskRunner(specStore, db)
	assert.NoError(t, err)

	// update the state of all tasks
	err = jobmine.JobStateWatcher(db)
	assert.NoError(t, err)

	// check the state of all jobs
	err = db.Where("run_id = ?", runId).First(&queryJob).Error
	assert.NoError(t, err)
	assert.Equal(t, queryJob.Status, jobmine.STATUS_SUCCESS)

	// check the state of tasks
	var queryTasks []jobmine.TaskRecord
	err = db.Where(&jobmine.TaskRecord{JobId: queryJob.ID}).Find(&queryTasks).Error
	assert.NoError(t, err)
	for _, queryTask := range queryTasks {
		assert.Equal(t, jobmine.STATUS_SUCCESS, queryTask.Status)
	}
}

func nukeConnections(db *gorm.DB) error {
	err := db.Delete(&data.Connection{}).Error
	if err != nil {
		return err
	}
	return db.Delete(&data.Mentorship{}).Error
}

func TestJobIntegration(t *testing.T) {
	theseTests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				users, err := createUsers(db)
				assert.NoError(t, err)

				runId := "seed_test_1"

				err = CreateSeedJob(db, runId, false,
					[]string{"COMPUTER_ENGINEERING", "SOFTWARE_ENGINEERING"},
					2021, 1, 100, nil, nil)
				assert.NoError(t, err)

				runAndTestRunners(t, db, runId)
				sqs := utility.QueueHelper.(utility.LocalQueueImpl)
				go sqs.QueueProcessor()

				var connections []data.Connection
				err = db.Model(&data.Connection{}).Preload("Mentorship").Find(&connections).Error
				assert.NoError(t, err)
				assert.Equal(t, 3, len(connections))

				matches := make([]match, 0)
				for _, connection := range connections {
					assert.NotNil(t, connection.Mentorship)
					assert.NotNil(t, connection.AcceptedAt)
					matches = append(matches, match{connection.UserOneId, connection.UserTwoId})
					assert.Equal(t, connection.Mentorship.MentorUserId, connection.UserOneId)
				}
				expectedMatches := []match{
					match{users[3].UserId, users[0].UserId},
					match{users[4].UserId, users[1].UserId},
					match{users[5].UserId, users[2].UserId},
				}
				assert.ElementsMatch(t, expectedMatches, matches)

				err = nukeConnections(db)
				assert.NoError(t, err)
			},
		},
		test.Test{
			Test: func(db *gorm.DB) {
				// Now testing with a restricted date range
				runId := "seed_test_2"
				now := time.Now()

				var users []data.User
				err := db.Order("user_id").Find(&users).Error
				assert.NoError(t, err)

				from := now.AddDate(0, 0, -1)
				to := now.AddDate(0, 0, 1)
				err = CreateSeedJob(db, runId, false,
					[]string{"COMPUTER_ENGINEERING", "SOFTWARE_ENGINEERING"},
					2021, 100, 1, &from, &to)
				assert.NoError(t, err)

				runAndTestRunners(t, db, runId)
				sqs := utility.QueueHelper.(utility.LocalQueueImpl)
				go sqs.QueueProcessor()

				var connections []data.Connection
				err = db.Model(&data.Connection{}).Preload("Mentorship").Find(&connections).Error
				assert.NoError(t, err)
				assert.Equal(t, 1, len(connections))

				// Only one match because lower years have at most one mentor, and only one upper year
				// was created within the date range.
				connection := connections[0]
				assert.NotNil(t, connection.Mentorship)
				assert.NotNil(t, connection.AcceptedAt)
				assert.Equal(t, users[4].UserId, connection.UserOneId)
				assert.Equal(t, users[1].UserId, connection.UserTwoId)
				assert.Equal(t, users[4].UserId, connection.Mentorship.MentorUserId)

				err = nukeConnections(db)
				assert.NoError(t, err)
			},
		},
		test.Test{
			Test: func(db *gorm.DB) {
				// Now testing with a dry run
				runId := "seed_test_3"

				err := CreateSeedJob(db, runId, true,
					[]string{"COMPUTER_ENGINEERING", "SOFTWARE_ENGINEERING"},
					2021, 100, 100, nil, nil)
				assert.NoError(t, err)

				runAndTestRunners(t, db, runId)

				var connections []data.Connection
				err = db.Model(&data.Connection{}).Preload("Mentorship").Find(&connections).Error
				assert.NoError(t, err)
				assert.Equal(t, 0, len(connections))
			},
		},
	}
	test.RunTestsWithDb(theseTests)
}
