package jobmine

import (
	"database/sql/driver"
	"encoding/json"
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
