package user

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/data"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
)

var testUserId = 1

func CreateUserForTest(t *testing.T, db *gorm.DB) *data.User {
	birthdate := "1996-10-07"
	req := api.SignupRequest{
		UserPersonalInfo: api.UserPersonalInfo{
			FirstName: "Firstname",
			LastName:  "Lastname",
			Gender:    0,
			Birthdate: &birthdate,
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
