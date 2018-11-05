package jobmine

import "database/sql/driver"

// TaskStatus Status of a particular job/task
type Status string

const (
	STATUS_CREATED Status = "CREATED"
	STATUS_RUNNING Status = "RUNNING"
	STATUS_SUCCESS Status = "SUCCESS"
	STATUS_FAILED  Status = "FAILED"
)

// Custom DB actions

func (u *Status) Scan(value interface{}) error { *u = Status(value.([]byte)); return nil }
func (u Status) Value() (driver.Value, error)  { return string(u), nil }
