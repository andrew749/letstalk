package query

import (
	"testing"

	"letstalk/server/core/test"
	"letstalk/server/core/utility/uw_email"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSearchUserFallback(t *testing.T) {
	test.RunTestWithDb(test.Test{
		Test: func(db *gorm.DB) {
			user, err := createTestUser(db, 1)
			assert.NoError(t, err)

			watEmail := uw_email.FromString("test@uwaterloo.ca")
			_, err = GenerateNewVerifyEmailId(db, 1, watEmail)
			assert.NoError(t, err)

			user2, err := GetUserByEmail(db, watEmail.ToStringRaw())
			assert.NoError(t, err)
			assert.Equal(t, user.FirstName, user2.FirstName)
			assert.Equal(t, user.Email, user2.Email)
		},
		TestName: "Test searching for user fallback",
	})
}
