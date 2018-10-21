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
	JobId uint `gorm:"primary_key"`

	// Denormalize to make it easier to debug this by looking at the db.
	JobType JobType `gorm:"primary_key"`

	// Running status of this task.
	Status Status

	// Metadata associated with this specific task.
	Metadata Metadata

	// Store blob error message
	ErrorData string `gorm:"type:text"`
}

func (r *TaskRecord) RecordSuccess(db *gorm.DB) error {
	return db.Where(r).Update("status = ?", Success).Error
}

func (r *TaskRecord) RecordRunning(db *gorm.DB) error {
	return db.Where(r).Update("status = ?", Running).Error
}

func (r *TaskRecord) RecordError(db *gorm.DB, err error) error {
	return db.Where(r).Update("status = ?", Failed).Update("error_data = ?", fmt.Sprintf("%+v", err)).Error
}

func (r *TaskRecord) GetJobRecordForTask(db *gorm.DB) (JobRecord, error) {
	var jobRecord JobRecord
	err := db.Where(&JobRecord{JobType: r.JobType}).First(&jobRecord).Error
	return jobRecord, err
}
