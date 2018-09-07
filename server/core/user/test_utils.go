package user

import (
	"testing"
	"github.com/jinzhu/gorm"
	"letstalk/server/data"
	"letstalk/server/core/api"
	"fmt"
	"github.com/stretchr/testify/require"
)

var testUserId = 1

func CreateUserForTest(t *testing.T, db *gorm.DB) *data.User {
	req := api.SignupRequest{
		UserPersonalInfo: api.UserPersonalInfo{
			FirstName: "Firstname",
			LastName:  "Lastname",
			Gender:    0,
			Birthdate: "1996-10-07",
		},
		Email:       fmt.Sprintf("test%d@test.com", testUserId),
		PhoneNumber: "5555555555",
		Password:    "test",
	}
	testUserId += 1
	usr, err := CreateUserWithAuth(db, req.Email, req.FirstName, req.LastName, req.Gender, req.Birthdate, data.USER_ROLE_DEFAULT, req.Password)
	require.NoError(t, err)
	// userId := c.Result.(struct{ UserId data.TUserID }).UserId
	return usr
}
