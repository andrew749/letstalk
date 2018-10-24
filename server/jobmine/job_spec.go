package jobmine

import (
	"github.com/jinzhu/gorm"
)

// GetTasksToCreate Function that creates db records for the instances of a job
// that should get run.
type GetTasksToCreate func(db *gorm.DB, jobRecord JobRecord) ([]*Metadata, error)

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
