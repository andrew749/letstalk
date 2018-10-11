package query

import (
	"letstalk/server/core/test"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSearchUserFallback(t *testing.T) {
	test.RunTestWithDb(test.Test{
		Test: func(db *gorm.DB) {
			user, err := createUser(db, 1)
			assert.NoError(t, err)

			watEmail := "test@uwaterloo.ca"
			_, err = GenerateNewVerifyEmailId(db, 1, watEmail)
			assert.NoError(t, err)

			user2, err := GetUserByEmail(db, watEmail)
			assert.NoError(t, err)
			assert.Equal(t, user.FirstName, user2.FirstName)
			assert.Equal(t, user.Email, user2.Email)
		},
		TestName: "Test searching for user fallback",
	})
}
