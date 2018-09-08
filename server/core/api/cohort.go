package api

import "letstalk/server/data"

// TODO, MOVE OFF THIS SOON!!!
type Cohort struct {
	CohortId   data.TCohortID `json:"cohortId"`
	ProgramId  string         `json:"programId"`
	SequenceId string         `json:"sequenceId"`
	GradYear   uint           `json:"gradYear"`
}

type CohortV2 struct {
	CohortId     data.TCohortID `json:"cohortId"`
	ProgramId    string         `json:"programId"`
	ProgramName  string         `json:"programName"`
	GradYear     uint           `json:"gradYear"`
	IsCoop       bool           `json:"isCoop"`
	SequenceId   *string        `json:"sequenceId"`
	SequenceName *string        `json:"sequenceName"`
}
