package jobmine

import (
	"database/sql/driver"

	"github.com/jinzhu/gorm"
)

// PopulateFunction Function that creates db records for the instances of a job
// that should get run.
type GetTasksToCreate func(db *gorm.DB, jobRecord JobRecord) ([]*Metadata, error)

// JobType Human readable identifier for a type of job. Used to determine which code is run.
type JobType string

// Custom DB actions
func (u *JobType) Scan(value interface{}) error { *u = JobType(value.(string)); return nil }
func (u JobType) Value() (driver.Value, error)  { return string(u), nil }

// JobSpec A definition of a job, essentially saying what code should get run.
// This also defines a function, `Populate` which is responsible for creating
// records for each instnatiation of a job.
type JobSpec struct {

	// JobType Unique identifier for the type of job that is more human readable.
	// Also how we index into the mapping of jobs to specs that should get run for a job.
	JobType JobType

	// Spec for each instance of a task that this job is composed of.
	TaskSpec TaskSpec

	// Populate Create task records to be executed at some point in the future
	GetTasksToCreate GetTasksToCreate
}
