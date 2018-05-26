package api

/**
 * Holds all the data that we currently associate with a user.
 */
type UserPersonalInfo struct {
	UserId    int    `json:"userId"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Gender    int    `json:"gender" binding:"required"`
	Birthdate int64  `json:"birthdate" binding:"required"` // unix time

	// TODO: Does this belong here?
	Secret     string  `json:"secret"`
	ProfilePic *string `json:"profilePic"`
}

type UserContactInfo struct {
	Email       string  `json:"email" binding:"required"`
	PhoneNumber *string `json:"phoneNumber"`
	FbId        *string `json:"fbId"`
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
