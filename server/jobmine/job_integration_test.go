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

const (
	testJobType  JobType = "Test Job"
	testJobType2 JobType = "Test Job 2"
)

func createTestJob(db *gorm.DB, runId string, jobType JobType) (*JobRecord, error) {
	return CreateJobRecord(
		db,
		runId,
		jobType,
		map[string]interface{}{},
		nil,
	)
}

func createTestTask(db *gorm.DB, jobId uint, runId string, jobType JobType) (*TaskRecord, error) {
	return CreateTaskRecord(
		db,
		jobId,
		runId,
		jobType,
		map[string]interface{}{},
	)
}

func resetDB(db *gorm.DB) {
	// HACK: we need a hard delete so that the primary key doesn't conflict with new inserted records.
	if err := db.Exec("DELETE FROM job_records;").Error; err != nil {
		panic(err)
	}

	if err := db.Exec("DELETE FROM task_records;").Error; err != nil {
		panic(err)
	}
}

func TestJobIntegration(t *testing.T) {
	test_light.RunTestsWithDb(
		provisionDb,
		[]test_light.Test{
			test_light.Test{
				TestName: "Schedule And run Job Success test",
				Test: func(db *gorm.DB) {
					resetDB(db)

					// create job
					const testRunId = "Test Run 1"
					now := time.Now()
					yesterday := now.AddDate(0, -1, 0) // yesterday

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
					job, err := CreateJobRecord(
						db,
						testRunId,
						testJobType,
						Metadata(map[string]interface{}{
							"job_key1": "job_value1",
						}),
						&yesterday,
					)
					assert.NoError(t, err)

					var queryJob JobRecord
					err = db.Where("run_id = ?", testRunId).Find(&queryJob).Error
					assert.NoError(t, err)
					assert.NotZero(t, queryJob.ID)
					assert.Equal(t, STATUS_CREATED, queryJob.Status)
					rlog.Info("Created job")

					// schedule tasks
					err = JobRunner(specStore, db, nil, nil)
					assert.NoError(t, err)

					// run tasks
					err = TaskRunner(specStore, db, nil, nil)
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
					resetDB(db)

					// DONE
				},
			},
		},
	)
}

const (
	runId1 = "job 1"
	runId2 = "job 2"
	runId3 = "job 3"
)

func TestGetJobmineJobToRun(t *testing.T) {
	test_light.RunTestsWithDb(
		provisionDb,
		[]test_light.Test{
			test_light.Test{
				TestName: "Test Getting normal jobmine jobs",
				Test: func(db *gorm.DB) {
					resetDB(db)
					rlog.Infof("test 1")
					job1, err := createTestJob(db, runId1, testJobType)
					assert.NoError(t, err)

					job2, err := createTestJob(db, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineJobsToRun(db, time.Now(), nil, nil)
					assert.NoError(t, err)
					assert.Len(t, jobs, 2)
					assert.Equal(t, jobs[0].ID, job1.ID)
					assert.Equal(t, jobs[1].ID, job2.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine jobs with runId filter",
				Test: func(db *gorm.DB) {
					resetDB(db)

					job1, err := createTestJob(db, runId1, testJobType)
					assert.NoError(t, err)

					_, err = createTestJob(db, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineJobsToRun(db, time.Now(), nil, []string{runId1})
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, job1.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine jobs with jobType filter",
				Test: func(db *gorm.DB) {
					resetDB(db)
					job1, err := createTestJob(db, runId1, testJobType)
					assert.NoError(t, err)

					_, err = createTestJob(db, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineJobsToRun(db, time.Now(), []JobType{testJobType}, nil)
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, job1.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine jobs with jobType filter and runId Filter",
				Test: func(db *gorm.DB) {
					resetDB(db)
					_, err := createTestJob(db, runId1, testJobType)
					assert.NoError(t, err)

					job2, err := createTestJob(db, runId2, testJobType2)
					assert.NoError(t, err)

					_, err = createTestJob(db, runId3, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineJobsToRun(db, time.Now(), []JobType{testJobType2}, []string{runId2})
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, job2.ID)
				},
			},
		})
}

const (
	jobId1 uint = 1
	jobId2 uint = 2
	jobId3 uint = 3
)

func TestGetJobmineTasksToRun(t *testing.T) {
	test_light.RunTestsWithDb(
		provisionDb,
		[]test_light.Test{
			test_light.Test{
				TestName: "Test Getting jobmine tasks",
				Test: func(db *gorm.DB) {
					resetDB(db)
					testTask, err := createTestTask(db, jobId1, runId1, testJobType)
					assert.NoError(t, err)

					testTask2, err := createTestTask(db, jobId2, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineTasksToRun(db, nil, nil)
					assert.NoError(t, err)
					assert.Len(t, jobs, 2)
					assert.Equal(t, jobs[0].ID, testTask.ID)
					assert.Equal(t, jobs[1].ID, testTask2.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine tasks by runId filter",
				Test: func(db *gorm.DB) {
					resetDB(db)
					testTask, err := createTestTask(db, jobId1, runId1, testJobType)
					assert.NoError(t, err)

					_, err = createTestTask(db, jobId2, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineTasksToRun(db, nil, []string{runId1})
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, testTask.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine tasks by jobType filter",
				Test: func(db *gorm.DB) {
					resetDB(db)
					testTask, err := createTestTask(db, jobId1, runId1, testJobType)
					assert.NoError(t, err)

					_, err = createTestTask(db, jobId2, runId2, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineTasksToRun(db, []JobType{testJobType}, nil)
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, testTask.ID)
				},
			},
			test_light.Test{
				TestName: "Test Getting jobmine tasks by runId and job type filter",
				Test: func(db *gorm.DB) {
					resetDB(db)
					_, err := createTestTask(db, jobId1, runId1, testJobType)
					assert.NoError(t, err)

					testTask2, err := createTestTask(db, jobId2, runId2, testJobType2)
					assert.NoError(t, err)

					_, err = createTestTask(db, jobId3, runId3, testJobType2)
					assert.NoError(t, err)

					jobs, err := GetJobmineTasksToRun(db, []JobType{testJobType2}, []string{runId2})
					assert.NoError(t, err)
					assert.Len(t, jobs, 1)
					assert.Equal(t, jobs[0].ID, testTask2.ID)
				},
			},
		})
}
