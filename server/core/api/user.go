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

type UserAdditionalData struct {
	MentorshipPreference *int    `json:"mentorshipPreference"`
	Bio                  *string `json:"bio"`
	Hometown             *string `json:"hometown"`
}

type UserContactInfo struct {
	Email       string  `json:"email" binding:"required"`
	PhoneNumber *string `json:"phoneNumber"`
	FbId        *string `json:"fbId"`
}

type MentorshipPreference uint

const (
	MENTORSHIP_PREFERENCE_MENTOR MentorshipPreference = iota + 1
	MENTORSHIP_PREFERENCE_MENTEE
)

type UserType int

// the roles a user can take in a relationship
const (
	USER_TYPE_MENTOR UserType = iota + 1
	USER_TYPE_MENTEE
	USER_TYPE_ASKER
	USER_TYPE_ANSWERER
	USER_TYPE_UNKNOWN UserType = -1
)

type UserVectorPreferenceType int

const (
	PREFERENCE_TYPE_ME UserVectorPreferenceType = iota
	PREFERENCE_TYPE_YOU
)
