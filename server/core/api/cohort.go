package api

import "letstalk/server/data"

type Cohort struct {
	CohortId   data.TCohortID `json:"cohortId"`
	ProgramId  string         `json:"programId"`
	SequenceId string         `json:"sequenceId"`
	GradYear   uint           `json:"gradYear"`
}
