package seed_mentorships_job

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_utility"
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
				users, err := test_helpers.CreateUsersForMatching(db)
				assert.NoError(t, err)

				runId := "seed_test_1"

				err = CreateSeedJob(db, runId, false,
					[]string{"COMPUTER_ENGINEERING", "SOFTWARE_ENGINEERING"},
					[]uint{2021, 2022},
					2021, 1, 100, nil, nil)
				assert.NoError(t, err)

				jobmine_utility.RunAndTestRunners(t, db, runId, specStore)
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
					[]uint{2021, 2022},
					2021, 100, 1, &from, &to)
				assert.NoError(t, err)

				jobmine_utility.RunAndTestRunners(t, db, runId, specStore)
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
					[]uint{2021, 2022},
					2021, 100, 100, nil, nil)
				assert.NoError(t, err)

				jobmine_utility.RunAndTestRunners(t, db, runId, specStore)

				var connections []data.Connection
				err = db.Model(&data.Connection{}).Preload("Mentorship").Find(&connections).Error
				assert.NoError(t, err)
				assert.Equal(t, 0, len(connections))
			},
		},
	}
	test.RunTestsWithDb(theseTests)
}
