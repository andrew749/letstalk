package data

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

type Cohort struct {
	CohortId   uint    `json:"cohortId" gorm:"not null;auto_increment;primary_key"`
	ProgramId  string  `json:"programId" gorm:"not null;unique_index:cohort_index"`
	Program    Program `gorm:"foreignkey:ProgramId;"`
	GradYear   uint    `json:"gradYear" gorm:"not null;unique_index:cohort_index"`
	SequenceId *string `json:"sequenceId" gorm:"not null;unique_index:cohort_index"`
}

// GetPrograms Return all valid programs
func GetPrograms() []string {
	return append(stream4Programs, stream8Programs...)
}

var (
	STREAM4 = "4STREAM"
	STREAM8 = "8STREAM"
)

var (
	stream8Programs = []string{
		"Software Engineering",
		"Electrical Engineering",
		"Computer Engineering",
		"Civil Engineering",
		"Management Engineering",
		"Nanotechnology Engineering",
		"Mechanical Engineering",
		"Mechatronics Engineering",
	}

	stream4Programs = []string{
		"Electrical Engineering",
		"Computer Engineering",
		"Environmental Engineering",
		"Geological Engineering",
		"Systems Design Engineering",
		"Mechanical Engineering",
		"Mechatronics Engineering",
	}
)

// GetSpecialCohorts Returns special cohorts that have sequence information.
func GetSpecialCohorts() []*Cohort {
	cohorts := make([]*Cohort, 0, 3*(2023-2018+1))

	for gradYear := 2018; gradYear <= 2023; gradYear++ {

		for _, program := range stream4Programs {
			cohorts = append(cohorts, &Cohort{
				ProgramId:  program,
				GradYear:   uint(gradYear),
				SequenceId: &STREAM4,
			})
		}

		for _, program := range stream8Programs {
			cohorts = append(cohorts, &Cohort{
				ProgramId:  program,
				GradYear:   uint(gradYear),
				SequenceId: &STREAM8,
			})
		}
	}

	return cohorts
}

// PopulateCohort Populate cohort db with info for special cohorts.
func PopulateCohort(db *gorm.DB) {
	cohorts := GetSpecialCohorts()

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

var reg, _ = regexp.Compile("[^a-zA-Z0-9_]+")

// normalize the program name by removing punctuation
func normalizeProgramName(programName string) string {
	return reg.ReplaceAllString(
		strings.Replace(strings.ToUpper(programName), " ", "_", -1), "")
}

// GetNormalizedProgramMapping Mapping of ProgramName Key -> Human Readable program name
func GetNormalizedProgramMapping() map[string]string {
	programs := GetPrograms()
	programMapping := make(map[string]string)

	for _, program := range programs {
		programMapping[normalizeProgramName(program)] = program
	}

	return programMapping
}

// GetNormalizedProgramMapping Mapping of ProgramName Key -> Human Readable program name
func GetReverseNormalizedProgramMapping() map[string]string {
	programs := GetPrograms()
	programMapping := make(map[string]string)

	for _, program := range programs {
		programMapping[program] = normalizeProgramName(program)
	}

	return programMapping
}

// GraduationYears Return all valid graduation years.
func ValidGraduationYears() []int {
	return []int{2018, 2019, 2020, 2021, 2022, 2023}
}
