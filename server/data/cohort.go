package data

import (
	"github.com/jinzhu/gorm"
)

type Cohort struct {
	CohortId  int     `json:"cohort_id" gorm:"not null;auto_increment;primary_key"`
	Program   Program `gorm:"foreignkey:ProgramId;"`
	ProgramId string  `json:"program_id" gorm:"not null;unique_index:cohort_index"`
	GradYear  int     `json:"grad_year" gorm:"not null;unique_index:cohort_index"`
	Sequence  string  `json:"sequence" gorm:"not null;unique_index:cohort_index"`
}

func PopulateCohort(db *gorm.DB) {
	cohorts := []*Cohort{
		&Cohort{
			ProgramId: "SOFTWARE_ENGINEERING",
			GradYear:  2019,
			Sequence:  "8STREAM",
		},
		&Cohort{
			ProgramId: "COMPUTER_ENGINEERING",
			GradYear:  2019,
			Sequence:  "8STREAM",
		},
		&Cohort{
			ProgramId: "COMPUTER_ENGINEERING",
			GradYear:  2019,
			Sequence:  "4STREAM",
		},
	}

	// add cohorts
	for _, cohort := range cohorts {
		db.FirstOrCreate(&cohort, Cohort{ProgramId: cohort.ProgramId, GradYear: cohort.GradYear, Sequence: cohort.Sequence})
	}
}
