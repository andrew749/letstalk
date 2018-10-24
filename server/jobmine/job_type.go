package jobmine

import "database/sql/driver"

// JobType Human readable identifier for a type of job. Used to determine which code is run.
type JobType string

// Custom DB actions

func (u *JobType) Scan(value interface{}) error { *u = JobType(value.(string)); return nil }
func (u JobType) Value() (driver.Value, error)  { return string(u), nil }
