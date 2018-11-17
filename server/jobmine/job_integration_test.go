package jobmine

import (
	"letstalk/server/core/test_light"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"

	"testing"

	"github.com/stretchr/testify/assert"
)

func provisionDb(db *gorm.DB) error {
	rlog.Info("Provisioning database for job_integration test.")
	if err := db.AutoMigrate(&JobRecord{}).Error; err != nil {
		rlog.Errorf("Could not create job record table: %+v", err)
		return err
	}
	if err := db.AutoMigrate(&TaskRecord{}).Error; err != nil {
		rlog.Errorf("Could not create task record table: %+v", err)
		return err
	}
	return nil
}

func TestJobIntegration(t *testing.T) {
	test_light.RunTestWithDb(
		provisionDb,
		test_light.Test{
			TestName: "Schedule And run Job Success test",
			Test: func(db *gorm.DB) {
				// create job
				const testJobType JobType = "Test Job"
				const testRunId = "Test Run 1"
				now := time.Now()
				job := JobRecord{
					Status:  STATUS_CREATED,
					JobType: testJobType,
					RunId:   testRunId,
					Metadata: Metadata(map[string]interface{}{
						"job_key1": "job_value1",
					}),
					StartTime: now.AddDate(0, -1, 0), // yesterday
				}
				// schedule job
				err := db.Create(&job).Error
				assert.NoError(t, err)

				var queryJob JobRecord
				err = db.Where("run_id = ?", testRunId).Find(&queryJob).Error
				assert.NoError(t, err)
				assert.NotZero(t, queryJob.ID)
				assert.Equal(t, STATUS_CREATED, queryJob.Status)
				rlog.Info("Created job")

				// define a test job
				specStore := JobSpecStore{
					JobSpecs: map[JobType]JobSpec{
						testJobType: JobSpec{
							JobType: testJobType,
							TaskSpec: TaskSpec{
								Execute: func(db *gorm.DB, jobRecord JobRecord, taskRecord TaskRecord) (interface{}, error) {
									assert.Equal(t, jobRecord.Status, STATUS_RUNNING)
									assert.Equal(t, jobRecord.JobType, testJobType)
									assert.Equal(t, jobRecord.RunId, testRunId)
									assert.Equal(t, jobRecord.Metadata["job_key1"], "job_value1")
									assert.Equal(t, taskRecord.JobId, jobRecord.ID)
									// HACK: the update times are slightly, off; this isnt an issue.
									now := time.Now()
									jobRecord.UpdatedAt = now
									taskRecord.Job.UpdatedAt = now
									assert.Equal(t, jobRecord, taskRecord.Job)
									assert.Equal(t, taskRecord.Status, STATUS_RUNNING)
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
							GetTasksToCreate: func(db *gorm.DB, jobRecord JobRecord) ([]Metadata, error) {
								metadata := Metadata(map[string]interface{}{
									"task_key1": "task_value1",
								})
								return []Metadata{metadata}, nil
							},
						},
					},
				}

				// schedule tasks
				err = JobRunner(specStore, db)
				assert.NoError(t, err)

				// run tasks
				err = TaskRunner(specStore, db)
				assert.NoError(t, err)

				// update the state of all tasks
				err = JobStateWatcher(db)
				assert.NoError(t, err)

				// check the state of all jobs
				err = db.Where("run_id = ?", queryJob.RunId).First(&queryJob).Error
				assert.NoError(t, err)
				assert.Equal(t, queryJob.Status, STATUS_SUCCESS)

				// check the state of tasks
				var queryTask TaskRecord
				err = db.Where(&TaskRecord{JobId: job.ID}).First(&queryTask).Error
				assert.NoError(t, err)
				assert.Equal(t, STATUS_SUCCESS, queryTask.Status)

				// DONE
			},
		},
	)
}
