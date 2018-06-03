package login

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/query"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	signupRequest := api.SignupRequest{
		UserPersonalInfo: api.UserPersonalInfo{
			FirstName: "Andrew",
			LastName:  "Codispoti",
			Gender:    0,
			Birthdate: 0,
		},
		Email:       "test@test.com",
		PhoneNumber: "5555555555",
		Password:    "test",
	}
	tests := []utility.Test{
		utility.Test{
			Test: func(db *gorm.DB) {
				var err error
				var user *data.User

				context := ctx.NewContext(nil, db, nil, nil)

				err = writeUser(&signupRequest, context)
				assert.NoError(t, err)

				userID := context.Result.(struct{ UserId int }).UserId
				user, err = query.GetUserById(db, userID)

				assert.NoError(t, err)
				assert.Equal(t, user.Email, signupRequest.Email)
			},
			TestName: "Test user creation",
		},
	}
	utility.RunTestsWithDb(tests)
}
