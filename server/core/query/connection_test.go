package query

import (
	"testing"

	"letstalk/server/core/test"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetAllMentorshipConnections(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := createTestUser(db, 1)
			assert.NoError(t, err)
			user2, err := createTestUser(db, 2)
			assert.NoError(t, err)
			user3, err := createTestUser(db, 3)
			assert.NoError(t, err)

			mentorship := data.Mentorship{
				MentorUserId: user1.UserId,
			}
			connection1 := data.Connection{
				UserOneId:  user1.UserId,
				UserTwoId:  user2.UserId,
				Mentorship: &mentorship,
			}
			err = db.Create(&connection1).Error
			assert.NoError(t, err)

			connection2 := data.Connection{
				UserOneId: user1.UserId,
				UserTwoId: user3.UserId,
			}
			err = db.Create(&connection2).Error
			assert.NoError(t, err)

			connections, err := GetAllMentorshipConnections(db)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(connections))
			assert.Equal(t, connection1.ConnectionId, connections[0].ConnectionId)
			assert.Equal(t, connection1.UserOneId, connections[0].UserOneId)
			assert.Equal(t, connection1.UserTwoId, connections[0].UserTwoId)
			assert.Equal(t, connection1.Mentorship.ConnectionId, connections[0].Mentorship.ConnectionId)
			assert.Equal(t, connection1.Mentorship.MentorUserId, connections[0].Mentorship.MentorUserId)
			assert.Equal(t, user1.Email, connections[0].UserOne.Email)
			assert.Equal(t, user2.Email, connections[0].UserTwo.Email)
		},
	}
	test.RunTestWithDb(thisTest)
}
