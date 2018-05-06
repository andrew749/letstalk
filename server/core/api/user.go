package api

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
	PREFERENCE_TYPE_ME UserVectorPreferenceType = iota
	PREFERENCE_TYPE_YOU
)
