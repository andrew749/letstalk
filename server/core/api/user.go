package api

/**
 * Holds all the data that we currently associate with a user.
 */
// TODO: Might want separate structs for /signup requests and /me responses, since they have
// different sets of required fields. It's better to have more structs with less optional fields
// than vice versa.
type User struct {
	// Personal information
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Gender    int    `json:"gender" binding:"required"`
	Birthday  int64  `json:"birthday" binding:"required"` // unix time

	// Contact information
	Email       string  `json:"email" binding:"required"`
	PhoneNumber *string `json:"phoneNumber"`
	FbId        *string `json:"fbId"`

	ProfilePic *string `json:"profilePic" binding:"required"`
}

type UserWithPassword struct {
	User
	Password string `json:"password"`
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
