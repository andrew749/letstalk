package api

import "letstalk/server/data"

type ProfileResponse struct {
	UserAdditionalData
	UserPersonalInfo
	UserContactInfo
	Cohort
	UserPositions    []UserPosition    `json:"userPositions"`
	UserSimpleTraits []UserSimpleTrait `json:"userSimpleTraits"`
}

type ProfileEditRequest struct {
	UserPersonalInfo
	UserAdditionalData
	PhoneNumber *string        `json:"phoneNumber"`
	CohortId    data.TCohortID `json:"cohortId" binding:"required"`
}

type ProfilePicResponse struct {
	ProfilePic *string `json:"profilePic"`
}
