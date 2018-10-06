package connection

import (
	"testing"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"

	"letstalk/server/core/sessions"
	"letstalk/server/core/user"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddMentorship(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				request := api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userTwo.Email,
				}
				requestError := handleAddMentorship(db, &request)
				assert.NoError(t, requestError)
				// Check all database tables are updated.
				conn, err := query.GetConnectionDetails(db, userOne.UserId, userTwo.UserId)
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
				request := api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userOne.Email,
				}
				// Same user id.
				assert.Error(t, handleAddMentorship(c.Db, &request))
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				connRequest := api.ConnectionRequest{
					UserId:        userTwo.UserId,
					IntentType:    data.INTENT_TYPE_SCAN_CODE,
				}
				handleRequestConnection(c, connRequest)
				// Connection already exists.
				request = api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userTwo.Email,
				}
				assert.Error(t, handleAddMentorship(c.Db, &request))
				request = api.CreateMentorshipByEmail{
					MentorEmail: "bademail@mail.com",
					MenteeEmail: userTwo.Email,
				}
				// No such user.
				assert.Error(t, handleAddMentorship(c.Db, &request))
			},
			TestName: "Test bad requests",
		},
	}
	test.RunTestsWithDb(tests)
}
