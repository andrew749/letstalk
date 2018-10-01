package connection

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"testing"

	"letstalk/server/core/sessions"
	"letstalk/server/core/user"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddMentorship(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				// Create new connection request
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				request := api.CreateMentorship{
					MentorId: userOne.UserId,
					MenteeId: userTwo.UserId,
				}
				requestError := handleAddMentorship(c, &request)
				assert.NoError(t, requestError)
				// Check all database tables are updated.
				conn, err := query.GetConnectionDetails(c.Db, userOne.UserId, userTwo.UserId)
				assert.NoError(t, err)
				assert.Equal(t, userOne.UserId, conn.UserOneId)
				assert.Equal(t, userTwo.UserId, conn.UserTwoId)
				assert.NotNil(t, conn.AcceptedAt)
				assert.Equal(t, data.INTENT_TYPE_ASSIGNED, conn.Intent.Type)
				assert.Equal(t, userOne.UserId, conn.Mentorship.MentorUserId)
			},
			TestName: "Test basic admin add mentorship",
		},
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				request := api.CreateMentorship{
					MentorId: userOne.UserId,
					MenteeId: userOne.UserId,
				}
				// Same user id.
				assert.Error(t, handleAddMentorship(c, &request))
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				connRequest := api.ConnectionRequest{
					UserId:        userTwo.UserId,
					IntentType:    data.INTENT_TYPE_SCAN_CODE,
				}
				handleRequestConnection(c, connRequest)
				// Connection already exists.
				request = api.CreateMentorship{
					MentorId: userOne.UserId,
					MenteeId: userTwo.UserId,
				}
				assert.Error(t, handleAddMentorship(c, &request))
				request = api.CreateMentorship{
					MentorId: 100,
					MenteeId: userTwo.UserId,
				}
				// No such user.
				assert.Error(t, handleAddMentorship(c, &request))
			},
			TestName: "Test bad requests",
		},
	}
	test.RunTestsWithDb(tests)
}
