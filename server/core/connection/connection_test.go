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

func TestRequestConnection(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				// Assert no connection exists at first.
				unexpected, err := query.GetConnectionDetailsUndirected(c.Db, userOne.UserId, userTwo.UserId)
				assert.NoError(t, err)
				assert.Nil(t, unexpected)
				// Create new connection request
				searchedTrait := "some trait"
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				request := api.ConnectionRequest{
					UserId:        userTwo.UserId,
					IntentType:    data.INTENT_TYPE_SEARCH,
					SearchedTrait: &searchedTrait,
				}
				result, err := HandleRequestConnection(c, request)
				assert.NoError(t, err)
				// Check result.
				expected := request
				expected.CreatedAt = result.CreatedAt
				assert.Equal(t, expected, *result)
				// Check data value from database.
				queried_for, err := query.GetConnectionDetailsUndirected(c.Db, userTwo.UserId, userOne.UserId)
				assert.Equal(t, userOne.UserId, queried_for.UserOneId)
				assert.Equal(t, userTwo.UserId, queried_for.UserTwoId)
				assert.Nil(t, queried_for.AcceptedAt)
				assert.Equal(t, request.IntentType, queried_for.Intent.Type)
				assert.Equal(t, *request.SearchedTrait, *queried_for.Intent.SearchedTrait)
				// Assert directed query still returns nil.
				unexpected, err = query.GetConnectionDetails(c.Db, userTwo.UserId, userOne.UserId)
				assert.NoError(t, err)
				assert.Nil(t, unexpected)
			},
			TestName: "Test user connection request",
		},
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				searchedTrait := "some trait"
				request := api.ConnectionRequest{
					UserId:        userTwo.UserId,
					IntentType:    data.INTENT_TYPE_SEARCH,
					SearchedTrait: &searchedTrait,
				}
				HandleRequestConnection(c, request)
				// Accept the request as user two.
				c.SessionData = &sessions.SessionData{UserId: userTwo.UserId}
				acceptReq := api.AcceptConnectionRequest{
					UserId: userOne.UserId,
				}
				result, err := HandleAcceptConnection(c, acceptReq)
				// Assert that AcceptedAt appears in result.
				assert.NoError(t, err)
				assert.NotNil(t, result.AcceptedAt)
				// Assert that AcceptedAt is not nil in db data.
				details, dbErr := query.GetConnectionDetails(c.Db, userOne.UserId, userTwo.UserId)
				assert.NoError(t, dbErr)
				assert.NotNil(t, details.AcceptedAt)
			},
			TestName: "Test accepting a user connection request",
		},
	}
	test.RunTestsWithDb(tests)
}

func TestRequestConnectionBadRequests(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				// Try to create connection request with nonexistent user.
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				searchedTrait := "some trait"
				request := api.ConnectionRequest{
					UserId:        100,
					IntentType:    data.INTENT_TYPE_SEARCH,
					SearchedTrait: &searchedTrait,
				}
				result, err := HandleRequestConnection(c, request)
				assert.Error(t, err)
				assert.Nil(t, result)
			},
			TestName: "Test request connection bad user id",
		},
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				userTwo := user.CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: userOne.UserId}
				searchedTrait := "some trait"
				request := api.ConnectionRequest{
					UserId:        userTwo.UserId,
					IntentType:    data.INTENT_TYPE_SEARCH,
					SearchedTrait: &searchedTrait,
				}
				HandleRequestConnection(c, request)
				// Accept the request as the same user.
				acceptReq := api.AcceptConnectionRequest{
					UserId: userTwo.UserId,
				}
				result, err := HandleAcceptConnection(c, acceptReq)
				// Assert that AcceptedAt appears in result.
				assert.Error(t, err)
				assert.Nil(t, result)
			},
			TestName: "Test accept connection no such request",
		},
	}
	test.RunTestsWithDb(tests)
}
