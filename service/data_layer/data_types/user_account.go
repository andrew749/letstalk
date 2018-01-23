package data_types

import (
	"log"
)

type UserAccount struct {
	Entity
	Name     string
	Email    string
	UserInfo UserInfo
}

func CreateUserAccount(name string, email string) {
	user := UserAccount{C, name, email}
	err := user.Save()

	if err != nil {
		log.Fatal("Unable to create user account")
	}
}

func (ua UserAccount) Save() error {
	// FIXME: save the database
}
