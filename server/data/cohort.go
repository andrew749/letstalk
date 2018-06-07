package data

import (
	"github.com/jinzhu/gorm"
)

type Cohort struct {
	CohortId   int     `json:"cohortId" gorm:"not null;auto_increment;primary_key"`
	ProgramId  string  `json:"programId" gorm:"not null;unique_index:cohort_index"`
	Program    Program `gorm:"foreignkey:ProgramId;"`
	GradYear   int     `json:"gradYear" gorm:"not null;unique_index:cohort_index"`
	SequenceId string  `json:"sequenceId" gorm:"not null;unique_index:cohort_index"`
}

func PopulateCohort(db *gorm.DB) {
	cohorts := make([]*Cohort, 0, 3*(2023-2018+1))

	for gradYear := 2018; gradYear <= 2023; gradYear++ {
		cohorts = append(cohorts, &Cohort{
			ProgramId:  "SOFTWARE_ENGINEERING",
			GradYear:   gradYear,
			SequenceId: "8STREAM",
		})
		cohorts = append(cohorts, &Cohort{
			ProgramId:  "COMPUTER_ENGINEERING",
			GradYear:   gradYear,
			SequenceId: "8STREAM",
		})
		cohorts = append(cohorts, &Cohort{
			ProgramId:  "COMPUTER_ENGINEERING",
			GradYear:   gradYear,
			SequenceId: "4STREAM",
		})
	}

	// add cohorts
	for _, cohort := range cohorts {
		db.FirstOrCreate(
			&cohort,
			Cohort{
				ProgramId:  cohort.ProgramId,
				GradYear:   cohort.GradYear,
				SequenceId: cohort.SequenceId,
			},
		)
	}
}
