package jobmine

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

// RunTask Runs a task record
// syncChannel is used to synchronize with the calling process since these are to be run in goroutines (no return code)
// taskRecord The actual task to run
func RunTask(db *gorm.DB, syncChannel chan<- TaskRecord, specStore JobSpecStore, taskRecord TaskRecord) error {

	// find the code to run
	taskSpec, err := specStore.GetTaskSpecForJobType(taskRecord.JobType)
	if err != nil {
		rlog.Errorf("Unable to find task spec for jobType=[%s]: %+v", taskRecord.JobType, err)
		syncChannel <- taskRecord
		taskRecord.RecordError(db, err)
		return err
	}

	// get the job metadata
	jobRecord, err := taskRecord.GetJobRecordForTask(db)
	if err != nil {
		rlog.Errorf("Unable to get jobRecord for jobType=[%s]: %+v", taskRecord.JobType, err)
		syncChannel <- taskRecord
		taskRecord.RecordError(db, err)
		return err
	}

	// Start running this job.
	// note the use of db rather than tx since we want all people to see this update
	taskRecord.RecordRunning(db)

	// Actually run the code
	tx := db.Begin()
	res, err := taskSpec.Execute(tx, jobRecord, taskRecord)

	// React to return value of code.
	if err != nil {
		rlog.Errorf("Task %d Failed: %+v\n", taskRecord.ID, err)
		tx.Rollback()

		// run failure callback
		taskSpec.OnError(tx, jobRecord, taskRecord, err)

		// write error status to job
		// note the use of db rather than tx since we want all people to see this update
		taskRecord.RecordError(db, err)

		// tell the runner that we're done
		syncChannel <- taskRecord

		// prevent bugs in case somebody writes code after the if statement
		return err
	} else {
		// write success status to job
		taskSpec.OnSuccess(tx, jobRecord, taskRecord, res)

		// write success status to job
		// note the use of db rather than tx since we want all people to see this update

		err = tx.Commit().Error
		if err != nil {
			rlog.Criticalf("Failed to commit changes to job.")
		}

		taskRecord.RecordSuccess(db)

		// tell the runner that we're done
		syncChannel <- taskRecord

		return nil
	}
}

// RunJob Creates task records for the job.
func RunJob(db *gorm.DB, specStore JobSpecStore, job JobRecord) error {
	tx := db.Begin()

	// get specs for job so we can get logic to create tasks
	rlog.Debugf("Fetching spec for jobType=[%s]", job.JobType)
	spec, err := specStore.GetJobSpecForJobType(job.JobType)
	if err != nil {
		tx.Rollback()
		job.SetJobStatus(db, STATUS_FAILED)
		rlog.Errorf("Unable to get spec for jobType=[%s]: %+v", job.JobType, err)
		return err
	}

	rlog.Debugf("Populating task records.")
	tasksMetadata, err := spec.GetTasksToCreate(tx, job)
	if err != nil {
		tx.Rollback()
		job.SetJobStatus(db, STATUS_FAILED)
		rlog.Errorf("Unable to get tasks metadata: %+v", err)
		return err
	}

	// create each of the tasks
	for _, taskMetadata := range tasksMetadata {
		if _, err := CreateTaskRecord(db, job.ID, job.RunId, job.JobType, taskMetadata); err != nil {
			tx.Rollback()
			job.SetJobStatus(db, STATUS_FAILED)
			rlog.Errorf("Unable to create task with metadata %+v", taskMetadata)
			return err
		}
		rlog.Infof("Successfully created task for jobId=[%d] jobType=[%s] runId=[%s] with metadata=[%+v]", job.ID, job.JobType, job.RunId, taskMetadata)
	}

	// update the job to running status
	if err := job.SetJobStatus(tx, STATUS_RUNNING); err != nil {
		tx.Rollback()
		job.SetJobStatus(db, STATUS_FAILED)
		return err
	}

	return tx.Commit().Error
}

// GetJobmineJobsToRun Jobs to run.
func GetJobmineJobsToRun(db *gorm.DB, latestJob time.Time, jobTypeWhitelist []JobType, runIdWhitelist []string) ([]JobRecord, error) {
	var jobs []JobRecord
	query := db.
		Where("status = ?", STATUS_CREATED).
		Where("start_time < ?", latestJob)

	if len(jobTypeWhitelist) > 0 {
		query = query.Where("job_type in (?)", jobTypeWhitelist)
	}

	if len(runIdWhitelist) > 0 {
		query = query.Where("run_id in (?)", runIdWhitelist)
	}

	// find all job records that are created but not started running
	if err := query.
		Find(&jobs).
		Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// JobRunner Finds Jobs that havent been executed yet and schedule them
// (by creating db records)
func JobRunner(jobSpecStore JobSpecStore, db *gorm.DB, jobTypeWhitelist []JobType, runIdWhitelist []string) error {
	jobs, err := GetJobmineJobsToRun(db, time.Now(), jobTypeWhitelist, runIdWhitelist)
	if err != nil {
		return err
	}

	rlog.Infof("Found %d jobs to run", len(jobs))
	for _, job := range jobs {
		rlog.Infof(
			"Running:\nJobType=[%s]\nRunId=[%s]\nMetadata=[%#v]",
			job.JobType,
			job.RunId,
			job.Metadata,
		)
		if err := RunJob(db, jobSpecStore, job); err != nil {
			rlog.Criticalf(
				"Error running:\nJobType=[%s]\nRunId=[%s]\nMetadata=[%#v]",
				job.JobType,
				job.RunId,
				job.Metadata,
			)
			return err
		}
	}
	return nil
}

// GetJobmineTasksToRun Tasks to run
func GetJobmineTasksToRun(db *gorm.DB, jobTypeWhitelist []JobType, runIdWhitelist []string) ([]TaskRecord, error) {
	var tasks []TaskRecord
	query := db.
		Where("status = ?", STATUS_CREATED).
		Preload("Job")

	if len(jobTypeWhitelist) > 0 {
		query = query.Where("job_type in (?)", jobTypeWhitelist)
	}

	if len(runIdWhitelist) > 0 {
		query = query.Where("run_id in (?)", runIdWhitelist)
	}

	if err := query.
		Find(&tasks).
		Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// TaskRunner Finds tasks that havent started yet and schedule them.
// whitelists are conjunctive as in they will filter the query
func TaskRunner(jobSpecStore JobSpecStore, db *gorm.DB, jobTypeWhitelist []JobType, runIdWhitelist []string) error {

	tasks, err := GetJobmineTasksToRun(db, jobTypeWhitelist, runIdWhitelist)
	if err != nil {
		return err
	}

	numTasks := len(tasks)
	// TODO(acod): change to a smaller number so we dont waste memory on larger number of tasks
	syncChannel := make(chan TaskRecord, numTasks)

	for _, task := range tasks {
		rlog.Infof("Running task %d: %#v", task.ID, task)
		// TODO(acod): Make these run in goroutines, need some retry logic since the transactions will
		// fail due to deadlocking issues.
		RunTask(db, syncChannel, jobSpecStore, task)
	}

	// wait for all tasks to finished, otherwise block
	for i := 0; i < numTasks; i++ {
		select {
		case taskComplete := <-syncChannel:
			rlog.Infof("Finished running task %d: %+v", taskComplete.ID, taskComplete)
		}
	}

	rlog.Infof("Finished running %d tasks", numTasks)

	return nil
}

// JobStateWatcher Watches the state of the job records and sees if a job is done, updating state
// To be run on a cron schedule
func JobStateWatcher(db *gorm.DB) error {
	var jobRecords []JobRecord
	err := db.Where("status = ?", STATUS_RUNNING).Find(&jobRecords).Error
	if err != nil {
		return err
	}

	// go over all running jobs
	for _, job := range jobRecords {
		taskRecords, err := GetTasksForJobId(db, job.ID)
		if err != nil {
			rlog.Warnf("Could not get tasks for jobId %d", job.ID)
			continue
		}
		var hasFailed = false
		var allComplete = true
		// go over all tasks for this job
		for _, task := range taskRecords {
			rlog.Debugf("Processing Task:\n\trunId=[%s]\n\ttaskId=[%d]\n", task.ID, task.RunId)
			if task.Status != STATUS_SUCCESS {
				rlog.Debugf("Found Incomplete Task:\n\trunId=[%s]\n\ttaskId=[%d]\n", task.ID, task.RunId)
				allComplete = false
			}
			if task.Status == STATUS_FAILED {
				hasFailed = true
				allComplete = false
				break
			}
		}

		// if a task failed, then this job failed
		if hasFailed {
			rlog.Debugf("Updating Job To failed:\n\tjobId=[%d]\n\trunId=[%s]", job.ID, job.RunId)
			err := job.SetJobStatus(db, STATUS_FAILED)
			if err != nil {
				rlog.Warnf("Unable to update status for job %d", job.ID)
			}
		}

		// if all tasks are complete, then this job is done successfully
		if allComplete {
			rlog.Debugf("Updating Job To Success:\n\tjobId=[%d]\n\trunId=[%s]", job.ID, job.RunId)
			err := job.SetJobStatus(db, STATUS_SUCCESS)
			if err != nil {
				rlog.Errorf("Unable to update job status: %+v", err)
			}
		}
	}
	return err
}
