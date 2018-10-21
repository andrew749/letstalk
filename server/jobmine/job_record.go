package jobmine

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// Metadata  Metadata for that is passed to each job and task
// run as part of the job.
type Metadata map[string]interface{}

// Custom DB actions

// Unmarshall map from json
func (u *Metadata) Scan(value interface{}) error {
	var tmp map[string]interface{}
	err := json.Unmarshal([]byte(value.(string)), &tmp)
	*u = tmp
	return err
}

// Marshall into json
func (u Metadata) Value() (driver.Value, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return string(data), err
}

// DATABASE RECORD
// JobRecord Instantiation of a specific job. A run of a job could have 0+ tasks
// each of which is a specific execution of the job.
type JobRecord struct {
	gorm.Model

	// JobName Human readable identifier for the type of job.
	// Used to determine which spec this job is running.
	JobType JobType `gorm:"primary_key"`

	// RunId Human readable unique identifier for the instantiation of this job.
	RunId string `gorm:"primary_key"`

	// Metadata that is part of this job that is accessible to each task running.
	// Allows us to configure and customize each specific job run.
	Metadata Metadata

	// When to start this job running. Allows us to run jobs "definitely after" a certain time.
	// This coupled with a frequent cron schedule makes it possible to get a bare bones job
	// scheduling system
	StartTime time.Time
}
