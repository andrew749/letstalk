package onboarding

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func CreateUserWithAuth(
	db *gorm.DB,
	email string,
	firstName string,
	lastName string,
	gender int,
	birthdate string,
	role data.UserRole,
	password string,
) (*data.User, error) {
	user, err := data.CreateUser(db, email, firstName, lastName, gender, birthdate, role)
	if err != nil {
		return nil, err
	}

	_, err = data.CreateAuthenticationData(db, user.UserId, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
