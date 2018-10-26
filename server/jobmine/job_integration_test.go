package jobmine

import (
	"letstalk/server/core/test"

	"github.com/jinzhu/gorm"

	"testing"

	"github.com/stretchr/testify/assert"
)

func JobIntegrationTest(t *testing.T) {
	test.RunTestWithDb(
		test.Test{
			TestName: "Schedule And run Job Success test",
			Test: func(db *gorm.DB) {
				// create job
				const testJobType JobType = "Test Job"
				const testRunId = "Test Run 1"
				job := JobRecord{
					Status:  Created,
					JobType: testJobType,
					RunId:   testRunId,
					Metadata: Metadata(map[string]interface{}{
						"job_key1": "job_value1",
					}),
				}
				// schedule job
				err := db.Create(&job).Error
				assert.NoError(t, err)

				// define a test job
				specStore := JobSpecStore{
					testJobType: JobSpec{
						JobType: testJobType,
						TaskSpec: TaskSpec{
							Execute: func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord) (interface{}, error) {
								assert.Equal(t, jobRecord.Status, Running)
								assert.Equal(t, jobRecord.JobType, testJobType)
								assert.Equal(t, jobRecord.RunId, testRunId)
								assert.Equal(t, jobRecord.Metadata["job_key1"], "job_value1")
								assert.Equal(t, taskRecord.JobId, jobRecord.ID)
								assert.Equal(t, taskRecord.JobRecord, jobRecord)
								assert.Equal(t, taskRecord.Status, Running)
								assert.Equal(t, taskRecord.Metadata["task_key1"], "task_value1")
								return "result_task1", nil
							},
							OnError: func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord, err error) {
								assert.Fail(t, "Should not have an error")
							},
							OnSuccess: func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord, res interface{}) {
								assert.Equal(t, res.(string), "result_task1")
							},
						},
						GetTasksToCreate: func(db *gorm.DB, jobRecord JobRecord) ([]*Metadata, error) {
							return []*Metadata{
								&map[string]interface{}{
									"task_key1": "task_value1",
								},
							}, nil
						},
					},
				}

				// schedule tasks
				err = JobRunner(db, specStore)
				assert.NoError(t, err)

				// run tasks
				err = TaskRunner(specStore, db)
				assert.NoError(t, err)

				// update the state of all tasks
				err = JobSpecWatcher(db)
				assert.NoError(t, err)

				// check the state of all jobs
				var queryJob JobRecord
				err = db.Where(&job).First(&queryJob).Error
				assert.NoError(t, err)
				assert.Equal(t, queryJob.Status, Success)

				// check the state of tasks
				var queryTask TaskRecord
				err = db.Where(&TaskRecord{JobId: job.ID}).First(&queryTask).Error
				assert.NoError(t, err)
				assert.Equal(t, queryTask.Status, Success)

				// DONE
			},
		},
	)
}
