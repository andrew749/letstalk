package api

type ProfileResponse struct {
	UserAdditionalData
	UserPersonalInfo
	UserContactInfo
	Cohort
}

type ProfileEditRequest struct {
	UserPersonalInfo
	UserAdditionalData
	PhoneNumber *string `json:"phoneNumber"`
	CohortId    int     `json:"cohortId" binding:"required"`
}

type ProfilePicResponse struct {
	ProfilePic *string `json:"profilePic"`
}
