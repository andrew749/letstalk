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
	RunId string `gorm:"size:190"`

	// DENORMALIZATION to make it easier to debug this by looking at the db.
	JobType JobType `gorm:"size:190"`

	// Running status of this task.
	Status Status `gorm:"not_null;size:50"`

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
	r.Status = STATUS_SUCCESS
	return db.Save(r).Error
}

// RecordRunning Mark the task as having started.
func (r *TaskRecord) RecordRunning(db *gorm.DB) error {
	r.Status = STATUS_RUNNING
	return db.Save(r).Error
}

// RecordError Mark the task as having failed.
func (r *TaskRecord) RecordError(db *gorm.DB, err error) error {
	r.Status = STATUS_FAILED
	r.ErrorData = fmt.Sprintf("%+v", err)
	return db.Save(r).Error
}

// GetJobRecordForTask Find a the job record associated with a task.
func (r *TaskRecord) GetJobRecordForTask(db *gorm.DB) (JobRecord, error) {
	var jobRecord JobRecord
	err := db.First(&jobRecord, r.JobId).Error
	return jobRecord, err
}

// CreateTaskRecord Creates a new task record in the db
func CreateTaskRecord(db *gorm.DB, jobId uint, runId string, jobType JobType, metadata Metadata) (*TaskRecord, error) {
	taskRecord := TaskRecord{
		JobId:    jobId,
		RunId:    runId,
		JobType:  jobType,
		Metadata: metadata,
		Status:   STATUS_CREATED,
	}
	if err := db.Save(&taskRecord).Error; err != nil {
		return nil, err
	}
	return &taskRecord, nil
}
