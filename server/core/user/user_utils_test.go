package user

import (
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserWithAuth(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				email := "andrew@test.com"
				firstName := "Andrew"
				lastName := "Codispoti"
				gender := 0
				birthdate := "1996-10-07"
				password := "test"

				user, err := CreateUserWithAuth(
					db,
					email,
					firstName,
					lastName,
					gender,
					birthdate,
					data.USER_ROLE_DEFAULT,
					password,
				)
				assert.NoError(t, err)
				assert.NotNil(t, user)
				user2, err := query.GetUserById(db, user.UserId)
				assert.NoError(t, err)
				assert.Equal(t, email, user2.Email)
				assert.Equal(t, firstName, user2.FirstName)
				assert.Equal(t, lastName, user2.LastName)
				assert.Equal(t, gender, user2.Gender)
				assert.Equal(t, birthdate, user2.Birthdate)
				hash, err := query.GetHashForUser(db, user.UserId)
				assert.NoError(t, err)
				assert.True(t, utility.CheckPasswordHash(password, *hash))
			},
			TestName: "Test user Creation with password",
		},
	}
	test.RunTestsWithDb(tests)
}
