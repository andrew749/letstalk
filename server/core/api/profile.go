package api

type Cohort struct {
	CohortId   int    `json:"cohortId" binding:"required"`
	ProgramId  string `json:"programId"`
	GradYear   int    `json:"gradYear"`
	SequenceId string `json:"sequenceId"`
}

type MyProfileResponse struct {
	UserPersonalInfo
	UserContactInfo
	Cohort
}

type ProfileEditRequest struct {
	UserPersonalInfo
	PhoneNumber *string `json:"phoneNumber"`
	CohortId    int     `json:"cohortId" binding:"required"`
}
