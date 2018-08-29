package user

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"testing"

	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"letstalk/server/core/sessions"
)

func createUserForTest(c *ctx.Context, t *testing.T) data.TUserID {
	signupRequest := api.SignupRequest{
		UserPersonalInfo: api.UserPersonalInfo{
			FirstName: "Andrew",
			LastName:  "Codispoti",
			Gender:    0,
			Birthdate: "1996-10-07",
		},
		Email:       "test@test.com",
		PhoneNumber: "5555555555",
		Password:    "test",
	}
	err := writeUser(&signupRequest, c)
	require.NoError(t, err)
	userId := c.Result.(struct{ UserId data.TUserID }).UserId
	return userId
}

func TestVerifyEmail(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				// Create user with initially unverified email.
				userId := createUserForTest(c, t)
				user, _ := query.GetUserById(db, userId)
				assert.False(t, user.IsEmailVerified)
				// Send request to send new account verification email.
				c.SessionData = &sessions.SessionData{UserId: userId}
				emailRequest := &api.SendAccountVerificationEmailRequest{Email: "foo@edu.uwaterloo.ca"}
				err := handleSendAccountVerificationEmailRequest(c, emailRequest)
				assert.NoError(t, err)
				// Look up verify email id.
				var verifyEmailId data.VerifyEmailId
				db.Where(&data.VerifyEmailId{UserId: userId}).First(&verifyEmailId)
				assert.Equal(t, userId, verifyEmailId.UserId)
				assert.False(t, verifyEmailId.IsUsed)
				assert.True(t, verifyEmailId.IsActive)
				assert.True(t, time.Now().Before(verifyEmailId.ExpirationDate))
				assert.Equal(t, "foo@edu.uwaterloo.ca", verifyEmailId.Email)
				// Send verification id to server (as if the verify link was clicked).
				verifyRequest := &api.VerifyEmailRequest{Id: verifyEmailId.Id}
				err = handleEmailVerification(c, verifyRequest)
				assert.NoError(t, err)
				// Assert that expected updates happened.
				var count uint
				db.Find(&data.VerifyEmailId{}, &data.VerifyEmailId{UserId: userId}).Count(&count)
				assert.Equal(t, uint(1), count)
				db.Where(&data.VerifyEmailId{UserId: userId}).First(&verifyEmailId)
				assert.False(t, verifyEmailId.IsActive)
				assert.True(t, verifyEmailId.IsUsed)
				user, _ = query.GetUserById(db, userId)
				assert.True(t, user.IsEmailVerified)
			},
			TestName: "Test user account email verification",
		},
	}
	test.RunTestsWithDb(tests)
}
