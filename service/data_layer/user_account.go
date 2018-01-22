package data_layer

import (
	"log"
)

type UserAccount struct {
	Entity
	Name  string
	Email string
}

func createNewUserAccount(name string, email string) {
	user := UserAccount{Entity{}, name, email}
	err := user.Save()

	if err != nil {
		log.Fatal("Unable to create user account")
	}
}
