package data_types

import (
	"log"
	"uwletstalk/service/data_layer"
)

type UserAccount struct {
	data_layer.Entity
	Name     string
	Email    string
	UserInfo *UserInfo
}

func CreateUserAccount(
	name string,
	email string,
	gender *Gender,
	birthdate *Birthdate,
) (*UserAccount, error) {
	userInfo, err := CreateUserInfo(gender, birthdate)

	if err != nil {
		return nil, err
	}
	user := UserAccount{
		data_layer.CreateEntity(),
		name,
		email,
		userInfo,
	}
	err = user.Save()

	if err != nil {
		log.Fatal("Unable to create user account")
		return nil, err
	}

	return &user, nil
}

func (ua UserAccount) Save() error {
	// FIXME: save the database
	return nil
}
