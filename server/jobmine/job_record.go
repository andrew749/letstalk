package jobmine

import (
	"time"

	"github.com/jinzhu/gorm"
)

// DATABASE RECORD

// JobRecord Instantiation of a specific job. A run of a job could have 0+ tasks
// each of which is a specific execution of the job.
type JobRecord struct {
	gorm.Model

	// The current status of the job
	Status Status

	// JobName Human readable identifier for the type of job.
	// Used to determine which spec this job is running.
	JobType JobType `gorm:"primary_key"`

	// RunId Human readable unique identifier for the instantiation of this job.
	// e.g. notifications_1996-10-07
	RunId string `gorm:"primary_key"`

	// Metadata that is part of this job that is accessible to each task running.
	// Allows us to configure and customize each specific job run.
	// Each task for a job has access to this data
	Metadata Metadata

	// When to start this job running. Allows us to run jobs "definitely after" a certain time.
	// This coupled with a frequent cron schedule makes it possible to get a bare bones job
	// scheduling system
	StartTime time.Time
}

// SetJobStatus Sets the status for a job
func (jr *JobRecord) SetJobStatus(db *gorm.DB, status Status) error {
	jr.Status = status
	return db.Save(jr).Error
}
