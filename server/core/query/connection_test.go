package query

import (
	"testing"
	"time"

	"letstalk/server/core/test"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestGetMentorshipConnectionsByDate(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := createTestUser(db, 1)
			assert.NoError(t, err)
			user2, err := createTestUser(db, 2)
			assert.NoError(t, err)
			user3, err := createTestUser(db, 3)
			assert.NoError(t, err)

			now := time.Now()
			connection1 := data.Connection{
				UserOneId: user1.UserId,
				UserTwoId: user2.UserId,
				Mentorship: &data.Mentorship{
					MentorUserId: user1.UserId,
					CreatedAt:    now,
				},
			}
			err = db.Create(&connection1).Error
			assert.NoError(t, err)

			connection2 := data.Connection{
				UserOneId: user2.UserId,
				UserTwoId: user3.UserId,
				Mentorship: &data.Mentorship{
					MentorUserId: user2.UserId,
					CreatedAt:    now.AddDate(0, 0, 10),
				},
			}
			err = db.Create(&connection2).Error
			assert.NoError(t, err)

			connection3 := data.Connection{
				UserOneId: user1.UserId,
				UserTwoId: user3.UserId,
				Mentorship: &data.Mentorship{
					MentorUserId: user1.UserId,
					CreatedAt:    now.AddDate(0, 0, -10),
				},
			}
			err = db.Create(&connection3).Error
			assert.NoError(t, err)

			startDate := now.AddDate(0, 0, -5)
			endDate := now.AddDate(0, 0, 5)
			connections, err := GetMentorshipConnectionsByDate(db, &startDate, &endDate)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(connections))
			assert.Equal(t, connection1.ConnectionId, connections[0].ConnectionId)
			assert.Equal(t, connection1.UserOneId, connections[0].UserOneId)
			assert.Equal(t, connection1.UserTwoId, connections[0].UserTwoId)
			assert.Equal(t, connection1.Mentorship.ConnectionId, connections[0].Mentorship.ConnectionId)
			assert.Equal(t, connection1.Mentorship.MentorUserId, connections[0].Mentorship.MentorUserId)
			assert.Equal(t, user1.Email, connections[0].UserOne.Email)
			assert.Equal(t, user2.Email, connections[0].UserTwo.Email)

			startDate = now.AddDate(0, 0, -10)
			endDate = now.AddDate(0, 0, 10)
			connections, err = GetMentorshipConnectionsByDate(db, &startDate, &endDate)
			assert.NoError(t, err)
			assert.Equal(t, 3, len(connections))

			startDate = now.AddDate(0, 0, 15)
			endDate = now.AddDate(0, 0, 20)
			connections, err = GetMentorshipConnectionsByDate(db, &startDate, &endDate)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(connections))
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetMentorshipConnectionsByDateNotSpecified(t *testing.T) {
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
				UserOneId:  user1.UserId,
				UserTwoId:  user3.UserId,
				Mentorship: nil,
			}
			err = db.Create(&connection2).Error
			assert.NoError(t, err)

			connections, err := GetMentorshipConnectionsByDate(db, nil, nil)
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
