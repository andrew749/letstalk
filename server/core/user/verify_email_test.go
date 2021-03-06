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
	"letstalk/server/core/sessions"
)

func TestVerifyEmail(t *testing.T) {
	theTest :=	test.Test{
		Test: func(db *gorm.DB) {
			c := ctx.NewContext(nil, db, nil, nil, nil)
			// Create user with initially unverified email.
			user := CreateUserForTest(t, c.Db)
			assert.False(t, user.IsEmailVerified)
			// Send request to send new account verification email.
			c.SessionData = &sessions.SessionData{UserId: user.UserId}
			emailRequest := &api.SendAccountVerificationEmailRequest{Email: "foo@edu.uwaterloo.ca"}
			err := handleSendAccountVerificationEmailRequest(c, emailRequest)
			assert.NoError(t, err)
			// Look up verify email id.
			var verifyEmailId data.VerifyEmailId
			db.Where(&data.VerifyEmailId{UserId: user.UserId}).First(&verifyEmailId)
			assert.Equal(t, user.UserId, verifyEmailId.UserId)
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
			db.Find(&data.VerifyEmailId{}, &data.VerifyEmailId{UserId: user.UserId}).Count(&count)
			assert.Equal(t, uint(1), count)
			var verifyEmailIdFinal data.VerifyEmailId
			db.Where(&data.VerifyEmailId{UserId: user.UserId}).First(&verifyEmailIdFinal)
			assert.False(t, verifyEmailIdFinal.IsActive)
			assert.True(t, verifyEmailIdFinal.IsUsed)
			assert.Equal(t, verifyEmailId.Id, verifyEmailIdFinal.Id)
			assert.Equal(t, verifyEmailId.UserId, verifyEmailIdFinal.UserId)
			assert.Equal(t, verifyEmailId.Email, verifyEmailIdFinal.Email)
			assert.Equal(t, verifyEmailId.ExpirationDate, verifyEmailIdFinal.ExpirationDate)
			user, _ = query.GetUserById(db, user.UserId)
			assert.True(t, user.IsEmailVerified)
		},
		TestName: "Test user account email verification",
	}
	test.RunTestWithDb(theTest)
}

func TestVerifyEmailMultipleRequests(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				// Create user with initially unverified email.
				user := CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: user.UserId}
				emailRequest := &api.SendAccountVerificationEmailRequest{Email: "foo@edu.uwaterloo.ca"}
				err := handleSendAccountVerificationEmailRequest(c, emailRequest)
				assert.NoError(t, err)
				// Generate a second account verification email.
				emailRequest = &api.SendAccountVerificationEmailRequest{Email: "bar@edu.uwaterloo.ca"}
				err = handleSendAccountVerificationEmailRequest(c, emailRequest)
				assert.NoError(t, err)
				// Look up verify email id entries.
				var verifyEmailIds []data.VerifyEmailId
				db.Where(&data.VerifyEmailId{UserId: user.UserId}).Order("expiration_date").Find(&verifyEmailIds)
				verifyEmailIdExpired, verifyEmailIdValid := verifyEmailIds[0], verifyEmailIds[1]
				assert.False(t, verifyEmailIdExpired.IsActive)
				assert.Equal(t, "foo@edu.uwaterloo.ca", verifyEmailIdExpired.Email)
				assert.True(t, verifyEmailIdValid.IsActive)
				assert.Equal(t, "bar@edu.uwaterloo.ca", verifyEmailIdValid.Email)
				// Send verification id of expired link.
				verifyRequest := &api.VerifyEmailRequest{Id: verifyEmailIdExpired.Id}
				err = handleEmailVerification(c, verifyRequest)
				assert.Error(t, err)
				// Send verification id of valid link.
				verifyRequest = &api.VerifyEmailRequest{Id: verifyEmailIdValid.Id}
				err = handleEmailVerification(c, verifyRequest)
				assert.NoError(t, err)
				// Assert that correct changes were made to table entries.
				db.Where(&data.VerifyEmailId{UserId: user.UserId}).Order("expiration_date").Find(&verifyEmailIds)
				verifyEmailIdExpired, verifyEmailIdValid = verifyEmailIds[0], verifyEmailIds[1]
				assert.False(t, verifyEmailIdExpired.IsActive)
				assert.False(t, verifyEmailIdValid.IsActive)
				assert.False(t, verifyEmailIdExpired.IsUsed)
				assert.True(t, verifyEmailIdValid.IsUsed)
			},
			TestName: "Test verify with multiple links",
		},
	}
	test.RunTestsWithDb(tests)
}

func TestVerifyEmailBadRequests(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				// Create user and verify their email.
				user := CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: user.UserId}
				emailRequest := &api.SendAccountVerificationEmailRequest{Email: "foo@edu.uwaterloo.ca"}
				handleSendAccountVerificationEmailRequest(c, emailRequest)
				var verifyEmailId data.VerifyEmailId
				db.Where(&data.VerifyEmailId{UserId: user.UserId}).First(&verifyEmailId)
				verifyRequest := &api.VerifyEmailRequest{Id: verifyEmailId.Id}
				handleEmailVerification(c, verifyRequest)
				// Try requesting another email verification for the user.
				err := handleSendAccountVerificationEmailRequest(c, emailRequest)
				assert.Error(t, err)
			},
			TestName: "Test verify email already verified",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				verifyRequest := &api.VerifyEmailRequest{Id: "BADID"}
				err := handleEmailVerification(c, verifyRequest)
				assert.Error(t, err)
			},
			TestName: "Test verify email invalid id",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				// Create user and verify their email.
				user := CreateUserForTest(t, c.Db)
				c.SessionData = &sessions.SessionData{UserId: user.UserId}
				emailRequest := &api.SendAccountVerificationEmailRequest{Email: "foo@edu.uwaterloo.ca"}
				handleSendAccountVerificationEmailRequest(c, emailRequest)
				var verifyEmailId data.VerifyEmailId
				db.Where(&data.VerifyEmailId{UserId: user.UserId}).First(&verifyEmailId)
				// Modify db to set expiration date to the past.
				verifyEmailId.ExpirationDate = time.Now().AddDate(0, 0, -1)
				db.Save(verifyEmailId)
				verifyRequest := &api.VerifyEmailRequest{Id: verifyEmailId.Id}
				// Try requesting email verification with the expired id.
				err := handleEmailVerification(c, verifyRequest)
				assert.Error(t, err)
			},
			TestName: "Test verify email expired link",
		},
	}
	test.RunTestsWithDb(tests)
}
