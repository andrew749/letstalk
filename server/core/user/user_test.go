package user

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"testing"

	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
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
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				var err error
				var tempUser *data.User

				context := ctx.NewContext(nil, db, nil, nil)

				err = writeUser(&signupRequest, context)
				assert.NoError(t, err)

				userID := context.Result.(struct{ UserId data.TUserID }).UserId
				tempUser, err = query.GetUserById(db, userID)

				assert.NoError(t, err)
				assert.Equal(t, tempUser.Email, signupRequest.Email)
			},
			TestName: "Test user creation",
		},
	}
	test.RunTestsWithDb(tests)
}

func TestBirthdate(t *testing.T) {
	type teststruct struct {
		msg       string
		birthdate string
		isValid   bool
	}
	tests := []teststruct{
		{"nominal", "1996-01-01", true},
		{"today", time.Now().Format(utility.BirthdateFormat), false},
		{"future", time.Now().AddDate(1, 0, 0).Format(utility.BirthdateFormat), false},
		{"edge valid", time.Now().AddDate(-13, 0, 0).Format(utility.BirthdateFormat), true},
		{"edge invalid", time.Now().AddDate(-13, 0, 1).Format(utility.BirthdateFormat), false},
	}
	for _, tempTest := range tests {
		isValid := validateUserBirthday(tempTest.birthdate) == nil
		assert.Equal(t, tempTest.isValid, isValid, "'%s' failed", tempTest.msg)
	}
}
