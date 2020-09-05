package test_job

import (
	"letstalk/server/jobmine"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

const TestJob jobmine.JobType = "TestJob"

var TestJobSpec jobmine.JobSpec = jobmine.JobSpec{
	JobType: TestJob,
	TaskSpec: jobmine.TaskSpec{
		Execute: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord) (interface{}, error) {
			rlog.Infof("Got data from taskRecord %s", taskRecord.Metadata["key"])
			return nil, nil
		},
		OnError: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, err error) {
			rlog.Errorf("Some error message")
		},
		OnSuccess: func(db *gorm.DB, jobRecord jobmine.JobRecord, taskRecord jobmine.TaskRecord, res interface{}) {
			rlog.Infof("Some success message")
		},
	},
	GetTasksToCreate: func(db *gorm.DB, jobRecord jobmine.JobRecord) ([]jobmine.Metadata, error) {
		res := make([]jobmine.Metadata, 0)
		data := jobmine.Metadata(map[string]interface{}{"key": "HELLO"})
		res = append(res, data)
		// Do some work
		return res, nil
	},
}

// CreateTestJob Creates a test job record
func CreateTestJob(db *gorm.DB, runId string, metadata jobmine.Metadata) error {
	_, err := jobmine.CreateJobRecord(db, runId, TestJob, metadata, nil)
	if err != nil {
		return err
	}
	return nil
}
