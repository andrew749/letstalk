package connection

import (
	"testing"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"

	"letstalk/server/core/sessions"
	"letstalk/server/core/user"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddMentorship(t *testing.T) {
	tests := []test.Test{
		{
			TestName: "Test basic admin add mentorship",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				request := api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userTwo.Email,
					RequestType: api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN,
				}
				requestError := HandleAddMentorship(db, &request)
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
		},
		{
			TestName: "Test bad requests",
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				request := api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userOne.Email,
					RequestType: api.CREATE_MENTORSHIP_TYPE_NOT_DRY_RUN,
				}
				// Same user id.
				assert.Error(t, HandleAddMentorship(c.Db, &request))
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				connRequest := api.ConnectionRequest{
					UserId:     userTwo.UserId,
					IntentType: data.INTENT_TYPE_SCAN_CODE,
				}
				HandleRequestConnection(c, connRequest)
				// Connection already exists.
				request = api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userTwo.Email,
				}
				assert.Error(t, HandleAddMentorship(c.Db, &request))
				request = api.CreateMentorshipByEmail{
					MentorEmail: "bademail@mail.com",
					MenteeEmail: userTwo.Email,
				}
				// No such user.
				assert.Error(t, HandleAddMentorship(c.Db, &request))
			},
		},
		{
			TestName: "Test dry run",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				request := api.CreateMentorshipByEmail{
					MentorEmail: userOne.Email,
					MenteeEmail: userTwo.Email,
					RequestType: api.CREATE_MENTORSHIP_TYPE_DRY_RUN,
				}
				requestError := HandleAddMentorship(db, &request)
				assert.NoError(t, requestError)
				// Check database tables are not updated.
				conn, err := query.GetConnectionDetails(db, userOne.UserId, userTwo.UserId)
				assert.NoError(t, err)
				assert.Nil(t, conn)
			},
		},
	}
	test.RunTestsWithDb(tests)
}

func checkAddMatchRoundMentorship(
	t *testing.T,
	db *gorm.DB,
	mentorUserId data.TUserID,
	menteeUserId data.TUserID,
	matchRoundId data.TMatchRoundID,
) {
	conn, err := query.GetMentorshipDetails(db, mentorUserId, menteeUserId)
	assert.NoError(t, err)
	assert.Equal(t, mentorUserId, conn.UserOneId)
	assert.Equal(t, menteeUserId, conn.UserTwoId)
	assert.NotNil(t, conn.AcceptedAt)
	assert.Equal(t, mentorUserId, conn.Mentorship.MentorUserId)
	assert.Equal(t, 1, len(conn.MatchRounds))
	assert.Equal(t, matchRoundId, conn.MatchRounds[0].MatchRoundId)
}

func TestAddMatchRoundMentorship(t *testing.T) {
	tests := []test.Test{
		{
			TestName: "Test full flow",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				matchRoundId := data.TMatchRoundID(10)

				err := AddMatchRoundMentorship(db, userOne.UserId, userTwo.UserId, matchRoundId)
				assert.NoError(t, err)

				checkAddMatchRoundMentorship(t, db, userOne.UserId, userTwo.UserId, matchRoundId)
			},
		},
		{
			TestName: "Test upgrade mentorship",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				matchRoundId := data.TMatchRoundID(10)

				test_helpers.CreateTestConnection(t, db, userOne.UserId, userTwo.UserId)
				err := AddMatchRoundMentorship(db, userOne.UserId, userTwo.UserId, matchRoundId)
				assert.NoError(t, err)

				checkAddMatchRoundMentorship(t, db, userOne.UserId, userTwo.UserId, matchRoundId)
			},
		},
		{
			TestName: "Test only write connection match round",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				matchRoundId := data.TMatchRoundID(10)

				connId := test_helpers.CreateTestConnection(t, db, userOne.UserId, userTwo.UserId)
				test_helpers.CreateTestMentorship(t, db, userOne.UserId, connId)
				err := AddMatchRoundMentorship(db, userOne.UserId, userTwo.UserId, matchRoundId)
				assert.NoError(t, err)

				checkAddMatchRoundMentorship(t, db, userOne.UserId, userTwo.UserId, matchRoundId)
			},
		},
		{
			TestName: "Test only write connection match round opposite dir",
			Test: func(db *gorm.DB) {
				// We don't switch the order of the connections when this happens
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				matchRoundId := data.TMatchRoundID(10)

				connId := test_helpers.CreateTestConnection(t, db, userTwo.UserId, userOne.UserId)
				test_helpers.CreateTestMentorship(t, db, userTwo.UserId, connId)
				err := AddMatchRoundMentorship(db, userOne.UserId, userTwo.UserId, matchRoundId)
				assert.NoError(t, err)

				checkAddMatchRoundMentorship(t, db, userTwo.UserId, userOne.UserId, matchRoundId)
			},
		},
		{
			TestName: "Test only write mentorship",
			Test: func(db *gorm.DB) {
				userOne := user.CreateUserForTest(t, db)
				userTwo := user.CreateUserForTest(t, db)
				matchRoundId := data.TMatchRoundID(10)

				connId := test_helpers.CreateTestConnection(t, db, userOne.UserId, userTwo.UserId)
				test_helpers.CreateTestConnectionMatchRound(t, db, connId, matchRoundId)
				err := AddMatchRoundMentorship(db, userOne.UserId, userTwo.UserId, matchRoundId)
				assert.NoError(t, err)

				checkAddMatchRoundMentorship(t, db, userOne.UserId, userTwo.UserId, matchRoundId)
			},
		},
	}
	test.RunTestsWithDb(tests)
}
