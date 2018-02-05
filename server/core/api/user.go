package api

import (
	"github.com/mijia/modelq/gmq"
)

/**
 * Holds all the data that we currently assiociate to a user.
 */
type User struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`

	// Contact information
	Email       string  `json:"email" binding:"required"`
	PhoneNumber *string `json:"phone_number"`

	// Personal information
	Gender   string `json:"gender" binding:"required"`
	Birthday int64  `json:"birthday" binding:"required"` // unix time

	Password *string `json:"password" binding:"required"`
}

/**
 * Try to see if there is school data assiociated with this account.
 * If there is no data, return nil
 */
func (u *User) GetSchoolInfo(db *gmq.Db) *SchoolInfo {
	// TODO(acod): get the school info for a particular user
	return nil
}
