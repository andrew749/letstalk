package bootstrap

import (
	"testing"

	"letstalk/server/core/api"
	"letstalk/server/core/connection"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
	"letstalk/server/utility"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func requestConnection(
	t *testing.T,
	db *gorm.DB,
	userOneId data.TUserID,
	userTwoId data.TUserID,
) {
	searchedTrait := "thing"
	message := "pls connect"
	request := api.ConnectionRequest{
		UserId:        userTwoId,
		IntentType:    data.INTENT_TYPE_SEARCH,
		SearchedTrait: &searchedTrait,
		Message:       &message,
	}
	_, err := connection.HandleRequestConnection(
		test_helpers.CreateTestContext(t, db, userOneId),
		request,
	)
	assert.NoError(t, err)
}

func acceptConnection(
	t *testing.T,
	db *gorm.DB,
	userOneId data.TUserID,
	userTwoId data.TUserID,
) {
	request := api.AcceptConnectionRequest{
		UserId: userTwoId,
	}
	_, err := connection.HandleAcceptConnection(
		test_helpers.CreateTestContext(t, db, userOneId),
		request,
	)
	assert.NoError(t, err)
}

func createMentorship(
	t *testing.T,
	db *gorm.DB,
	userOneId data.TUserID,
	userOneEmail string,
	userTwoEmail string,
) {
	request := api.CreateMentorshipByEmail{
		userOneEmail,
		userTwoEmail,
		api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN,
	}
	err := connection.HandleAddMentorship(db, &request)
	assert.NoError(t, err)
}

// This is actually a pretty big integration test of the system.
func TestGetCurrentUserBootstrapStatusController(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			users := make([]data.User, 0)
			for i := 1; i <= 7; i++ {
				user, err := test_helpers.CreateTestSetupUser(db, i)
				assert.NoError(t, err)
				users = append(users, *user)
			}

			requestConnection(t, db, users[0].UserId, users[1].UserId)
			requestConnection(t, db, users[2].UserId, users[0].UserId)
			requestConnection(t, db, users[0].UserId, users[3].UserId)
			acceptConnection(t, db, users[3].UserId, users[0].UserId)
			requestConnection(t, db, users[4].UserId, users[0].UserId)
			acceptConnection(t, db, users[0].UserId, users[4].UserId)
			createMentorship(t, db, users[0].UserId, users[0].Email, users[5].Email)
			createMentorship(t, db, users[6].UserId, users[6].Email, users[0].Email)

			// Need this to purge the queue
			sqs := utility.QueueHelper.(utility.LocalQueueImpl)
			go sqs.QueueProcessor()

			c := test_helpers.CreateTestContext(t, db, users[0].UserId)
			err := GetCurrentUserBoostrapStatusController(c)
			assert.NoError(t, err)

			res := c.Result.(api.BootstrapResponse)
			assert.Equal(t, api.ACCOUNT_MATCHED, res.State)
			assert.Equal(t, 1, len(res.Connections.OutgoingRequests))
			assert.Equal(t, users[1].UserId, res.Connections.OutgoingRequests[0].UserId)
			assert.Equal(t, 1, len(res.Connections.IncomingRequests))
			assert.Equal(t, users[2].UserId, res.Connections.IncomingRequests[0].UserId)
			peers := res.Connections.Peers
			assert.Equal(t, 2, len(peers))
			peerUserIds := map[data.TUserID]interface{}{
				peers[0].UserProfile.UserId: nil,
				peers[1].UserProfile.UserId: nil,
			}
			_, hasOne := peerUserIds[users[3].UserId]
			_, hasTwo := peerUserIds[users[4].UserId]
			assert.True(t, hasOne)
			assert.True(t, hasTwo)
			assert.Equal(t, 1, len(res.Connections.Mentors))
			assert.Equal(t, users[5].UserId, res.Connections.Mentors[0].UserProfile.UserId)
			assert.Equal(t, 1, len(res.Connections.Mentees))
			assert.Equal(t, users[6].UserId, res.Connections.Mentees[0].UserProfile.UserId)
		},
	}
	test.RunTestWithDb(thisTest)
}

// Test that we short-circuit if not done onboarding
func TestGetCurrentUserBootstrapStatusControllerNotFinishedOnboarding(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			users := make([]data.User, 0)
			for i := 1; i <= 7; i++ {
				user, err := test_helpers.CreateTestSetupUser(db, i)
				assert.NoError(t, err)
				users = append(users, *user)
			}
			users[0].IsEmailVerified = false
			dbErr := db.Save(&users[0]).Error
			assert.NoError(t, dbErr)

			// Need this to purge the queue
			sqs := utility.QueueHelper.(utility.LocalQueueImpl)
			go sqs.QueueProcessor()

			requestConnection(t, db, users[0].UserId, users[1].UserId)
			requestConnection(t, db, users[2].UserId, users[0].UserId)
			requestConnection(t, db, users[0].UserId, users[3].UserId)
			acceptConnection(t, db, users[3].UserId, users[0].UserId)
			requestConnection(t, db, users[4].UserId, users[0].UserId)
			acceptConnection(t, db, users[0].UserId, users[4].UserId)
			createMentorship(t, db, users[0].UserId, users[0].Email, users[5].Email)
			createMentorship(t, db, users[6].UserId, users[6].Email, users[0].Email)

			c := test_helpers.CreateTestContext(t, db, users[0].UserId)
			err := GetCurrentUserBoostrapStatusController(c)
			assert.NoError(t, err)

			res := c.Result.(api.BootstrapResponse)
			assert.Equal(t, api.ACCOUNT_CREATED, res.State)
			assert.Equal(t, 0, len(res.Connections.OutgoingRequests))
			assert.Equal(t, 0, len(res.Connections.IncomingRequests))
			assert.Equal(t, 0, len(res.Connections.Peers))
			assert.Equal(t, 0, len(res.Connections.Mentors))
			assert.Equal(t, 0, len(res.Connections.Mentees))
		},
	}
	test.RunTestWithDb(thisTest)
}
