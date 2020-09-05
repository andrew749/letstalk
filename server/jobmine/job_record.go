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
	Status Status `gorm:"not_null;size:50"`

	// JobName Human readable identifier for the type of job.
	// Used to determine which spec this job is running.
	JobType JobType `gorm:"primary_key;size:190"`

	// RunId Human readable unique identifier for the instantiation of this job.
	// e.g. notifications_1996-10-07
	RunId string `gorm:"unique;size:190"`

	// Metadata that is part of this job that is accessible to each task running.
	// Allows us to configure and customize each specific job run.
	// Each task for a job has access to this data
	Metadata Metadata `gorm:"type:text"`

	// When to start this job running. Allows us to run jobs "definitely after" a certain time.
	// This coupled with a frequent cron schedule makes it possible to get a bare bones job
	// scheduling system
	StartTime time.Time `gorm:"" sql:"DEFAULT:current_timestamp"`
}

// SetJobStatus Sets the status for a job
func (jr *JobRecord) SetJobStatus(db *gorm.DB, status Status) error {
	jr.Status = status
	return db.Save(jr).Error
}

func CreateJobRecord(db *gorm.DB, runId string, jobType JobType, metadata Metadata, startTime *time.Time) (*JobRecord, error) {

	// by default start now
	if startTime == nil {
		var now = time.Now()
		startTime = &now
	}

	jobRecord := JobRecord{
		RunId:     runId,
		JobType:   jobType,
		Metadata:  metadata,
		Status:    STATUS_CREATED,
		StartTime: *startTime,
	}
	if err := db.Save(&jobRecord).Error; err != nil {
		return nil, err
	}
	return &jobRecord, nil
}
