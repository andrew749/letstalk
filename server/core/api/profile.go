package api

type Cohort struct {
	CohortId   int    `json:"cohortId" binding:"required"`
	ProgramId  string `json:"programId"`
	GradYear   int    `json:"gradYear"`
	SequenceId string `json:"sequenceId"`
}

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
