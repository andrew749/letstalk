package api

import "letstalk/server/data"

type RelationshipType string

const (
	RELATIONSHIP_TYPE_NONE           RelationshipType = "NONE"
	RELATIONSHIP_TYPE_YOU_REQUESTED  RelationshipType = "YOU_REQUESTED"
	RELATIONSHIP_TYPE_THEY_REQUESTED RelationshipType = "THEY_REQUESTED"
	RELATIONSHIP_TYPE_CONNECTED      RelationshipType = "CONNECTED"
)

type ProfileResponse struct {
	UserAdditionalData
	UserPersonalInfo
	UserContactInfo
	Cohort
	UserPositions    []UserPosition    `json:"userPositions"`
	UserSimpleTraits []UserSimpleTrait `json:"userSimpleTraits"`
}

type MatchProfileResponse struct {
	ProfileResponse
	RelationshipType RelationshipType `json:"relationshipType"`
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
