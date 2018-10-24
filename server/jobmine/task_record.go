package jobmine

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// DATABASE RECORD

// TaskRecord An instantiation of a specific task. This keeps track of state for a task.
type TaskRecord struct {
	gorm.Model

	// JobId The actual job this task is part of.
	Job   JobRecord `gorm:"foreignkey:JobId;association_foreignkey:ID"`
	JobId uint      `gorm:"primary_key;auto_increment:false"`

	// DENORMALIZATION to make it easier to look for these things in the db.
	RunId string

	// DENORMALIZATION to make it easier to debug this by looking at the db.
	JobType JobType

	// Running status of this task.
	Status Status

	// Metadata associated with this specific task.
	Metadata Metadata `gorm:"type:text"`

	// Store blob error message
	ErrorData string `gorm:"type:text"`
}

// GetTasksForJobId Find all tasks corresponding to a job id.
func GetTasksForJobId(db *gorm.DB, jobId uint) ([]*TaskRecord, error) {
	var records []*TaskRecord
	err := db.Where(&TaskRecord{JobId: jobId}).Find(&records).Error
	return records, err
}

// RecordSuccess Mark the task as having completed successfully.
func (r *TaskRecord) RecordSuccess(db *gorm.DB) error {
	return db.Where(r).Update("status = ?", Success).Error
}

// RecordRunning Mark the task as having started.
func (r *TaskRecord) RecordRunning(db *gorm.DB) error {
	return db.Where(r).Update("status = ?", Running).Error
}

// RecordError Mark the task as having failed.
func (r *TaskRecord) RecordError(db *gorm.DB, err error) error {
	return db.Where(r).Update("status = ?", Failed).Update("error_data = ?", fmt.Sprintf("%+v", err)).Error
}

// GetJobRecordForTask Find a the job record associated with a task.
func (r *TaskRecord) GetJobRecordForTask(db *gorm.DB) (JobRecord, error) {
	var jobRecord JobRecord
	err := db.Where(&JobRecord{JobType: r.JobType}).First(&jobRecord).Error
	return jobRecord, err
}
