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
