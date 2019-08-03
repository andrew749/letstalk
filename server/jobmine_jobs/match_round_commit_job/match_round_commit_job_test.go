package match_round_commit_job

import (
	"fmt"
	"testing"

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
		MATCH_ROUND_COMMIT_JOB: CommitJobSpec,
	},
}

func TestJobIntegration(t *testing.T) {
	theseTests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				users, err := test_helpers.CreateUsersForMatching(db)
				assert.NoError(t, err)

				matchRound, err := test_helpers.CreateTestMatchRound(db)
				expectedMatches := []match{
					match{users[3].UserId, users[0].UserId},
					match{users[4].UserId, users[1].UserId},
					match{users[5].UserId, users[2].UserId},
					match{users[6].UserId, users[7].UserId},
				}
				for _, match := range expectedMatches {
					matchRoundMatch := data.MatchRoundMatch{
						MatchRoundId: matchRound.Id,
						MenteeUserId: match.menteeId,
						MentorUserId: match.mentorId,
						Score:        123,
					}
					err = db.Save(&matchRoundMatch).Error
					assert.NoError(t, err)
				}

				// Create some existing connections to see if idempotency works
				test_helpers.CreateTestConnection(t, db, users[4].UserId, users[1].UserId)
				connId1 := test_helpers.CreateTestConnection(t, db, users[5].UserId, users[2].UserId)
				connId2 := test_helpers.CreateTestConnection(t, db, users[6].UserId, users[7].UserId)
				test_helpers.CreateTestMentorship(t, db, users[5].UserId, connId1)
				test_helpers.CreateTestConnectionMatchRound(t, db, connId2, matchRound.Id)

				runId, err := CreateCommitJob(db, matchRound.Id)
				assert.NoError(t, err)
				assert.Equal(t, fmt.Sprintf("match-round-commit-%d", matchRound.Id), *runId)

				jobmine_utility.RunAndTestRunners(t, db, *runId, specStore)
				sqs := utility.QueueHelper.(utility.LocalQueueImpl)
				go sqs.QueueProcessor()

				var connections []data.Connection
				err = db.Model(
					&data.Connection{},
				).Preload("Mentorship").Preload("MatchRounds").Find(&connections).Error
				assert.NoError(t, err)
				assert.Equal(t, len(expectedMatches), len(connections))

				matches := make([]match, 0)
				for _, connection := range connections {
					assert.NotNil(t, connection.Mentorship)
					assert.NotNil(t, connection.AcceptedAt)
					matches = append(matches, match{connection.UserOneId, connection.UserTwoId})
					assert.Equal(t, connection.Mentorship.MentorUserId, connection.UserOneId)
					assert.Equal(t, 1, len(connection.MatchRounds))
					assert.Equal(t, matchRound.Id, connection.MatchRounds[0].MatchRoundId)
				}
				assert.ElementsMatch(t, expectedMatches, matches)
			},
		},
	}
	test.RunTestsWithDb(theseTests)
}
