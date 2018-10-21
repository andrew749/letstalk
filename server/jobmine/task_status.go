package jobmine

import "database/sql/driver"

// TaskStatus Status of a particular job
type Status string

const (
	Created Status = "CREATED"
	Running        = "RUNNING"
	Success        = "SUCCESS"
	Failed         = "FAILED"
)

// Custom DB actions
func (u *Status) Scan(value interface{}) error { *u = Status(value.(string)); return nil }
func (u Status) Value() (driver.Value, error)  { return string(u), nil }