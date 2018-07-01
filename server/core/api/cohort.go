package api

type Cohort struct {
	CohortId   uint    `json:"cohortId"`
	ProgramId  string  `json:"programId"`
	SequenceId *string `json:"sequenceId"`
	GradYear   uint    `json:"gradYear"`
}
