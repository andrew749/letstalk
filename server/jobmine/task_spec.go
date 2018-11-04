package jobmine

import "github.com/jinzhu/gorm"

// TaskSpec Information needed to run a specific task.
// Essentially which functions define a task.
type TaskSpec struct {

	// Execute Code to run for a specific task
	Execute func(
		db *gorm.DB,
		jobRecord JobRecord, // To get job specific configuration infromation
		taskRecord TaskRecord, // To get task specific configuration information
	) (interface{}, error)

	// Callbacks based on job status
	OnError   func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord, err error)
	OnSuccess func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord, res interface{})
}
