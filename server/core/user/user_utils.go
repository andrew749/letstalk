package user

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// CreateUserWithAuth creates a user with the following password
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
	tempUser, err := data.CreateUser(db, email, firstName, lastName, gender, birthdate, role)
	if err != nil {
		return nil, err
	}

	_, err = data.CreateAuthenticationData(db, tempUser.UserId, password)
	if err != nil {
		return nil, err
	}
	return tempUser, nil
}
