package api

import (
	"errors"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
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

type UserType string

// the roles a user can take in a relationship
const (
	USER_TYPE_MENTOR  UserType = "user_type_mentor"
	USER_TYPE_MENTEE  UserType = "user_type_mentee"
	USER_TYPE_UNKNOWN UserType = "user_type_unknown"
)

type UserVectorPreferenceType int

const (
	MePreference UserVectorPreferenceType = iota
	YouPreference
)

func GetUserWithId(db *gorm.DB, userId int) (*data.User, error) {
	var user data.User
	if db.Where("user_id = ?", userId).First(&user).RecordNotFound() {
		return nil, errors.New("Unable to find user")
	}

	return &user, nil
}

func GetFullUserWithId(db *gorm.DB, userId int) (*data.User, error) {
	var user data.User
	if err := db.Where("user_id = ?", userId).
		Preload("Mentees").
		Preload("Mentors").First(&user).Error; err != nil {
		rlog.Debug(err.Error())
		return nil, err
	}
	return &user, nil
}
