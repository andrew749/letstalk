package api

import (
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

/**
 * Holds all the data that we currently associate with a user.
 */
type User struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`

	// Contact information
	Email       string  `json:"email" binding:"required"`
	PhoneNumber *string `json:"phoneNumber"`

	// Personal information
	Gender   string `json:"gender" binding:"required"`
	Birthday int64  `json:"birthday" binding:"required"` // unix time

	Password *string `json:"password" binding:"required"`
}

func GetUserWithId(userId int, db *gmq.Db) (*data.User, error) {
	userObj := data.UserObjs
	user, err := userObj.Select().Where(userObj.FilterUserId("=", userId)).One(db)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
