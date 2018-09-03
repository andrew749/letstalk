package converters

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

func ApiCohortV2FromDataCohort(cohort *data.Cohort) *api.CohortV2 {
	return &api.CohortV2{
		CohortId:     cohort.CohortId,
		ProgramId:    cohort.ProgramId,
		ProgramName:  cohort.ProgramName,
		IsCoop:       cohort.IsCoop,
		GradYear:     cohort.GradYear,
		SequenceId:   cohort.SequenceId,
		SequenceName: cohort.SequenceName,
	}
}
