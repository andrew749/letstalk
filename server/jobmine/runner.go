package jobmine

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

func RunTask(db *gorm.DB, syncChannel chan TaskRecord, specStore JobSpecStore, taskRecord TaskRecord) {
	tx := db.Begin()
	taskSpec, err := specStore.GetTaskSpecForJobType(taskRecord.JobType)
	if err != nil {
		rlog.Errorf("Unable to find task spec for jobType=[%s]: %+v", taskRecord.JobType, err)
		syncChannel <- taskRecord
		return
	}
	jobRecord, err := taskRecord.GetJobRecordForTask(db)
	if err != nil {
		rlog.Errorf("Unable to get jobRecord for jobType=[%s]: %+v", taskRecord.JobType, err)
		syncChannel <- taskRecord
		return
	}

	taskRecord.RecordRunning(tx)
	res, err := taskSpec.Execute(tx, jobRecord, taskRecord)
	if err != nil {
		rlog.Errorf("Task Failed:\n")
		tx.Rollback()
		taskSpec.OnError(tx, jobRecord, taskRecord, err)
		// write error status to job
		taskRecord.RecordError(tx, err)

		// tell the runner that we're done
		syncChannel <- taskRecord
		return
	} else {
		// write success status to job
		taskSpec.OnSuccess(tx, jobRecord, taskRecord, res)
		taskRecord.RecordSuccess(tx)
		err = tx.Commit().Error
		if err != nil {
			rlog.Criticalf("Failed to commit changes to job.")
		}

		// tell the runner that we're done
		syncChannel <- taskRecord
		return
	}
}

func RunJob(db *gorm.DB, specStore JobSpecStore, job JobRecord) error {
	tx := db.Begin()
	rlog.Debugf("Fetching spec for jobType=[%s]", job.JobType)
	spec, err := specStore.GetJobSpecForJobtype(job.JobType)
	if err != nil {
		tx.Rollback()
		rlog.Errorf("Unable to get spec for jobType=[%s]: %+v", job.JobType, err)
		return err
	}

	rlog.Debugf("Populating task records.")
	tasksMetadata, err := spec.GetTasksToCreate(tx, job)
	if err != nil {
		tx.Rollback()
		rlog.Errorf("Unable to get tasks metadata: %+v", err)
		return err
	}

	for _, taskMetadata := range tasksMetadata {
		if err := tx.Create(&TaskRecord{
			JobId:    job.ID,
			JobType:  job.JobType,
			Status:   Created,
			Metadata: *taskMetadata,
		}).Error; err != nil {
			tx.Rollback()
			rlog.Errorf("Unable to create task with metadata %+v", taskMetadata)
			return err
		}
		rlog.Infof("Successfully created task for jobId=[%d] jobType=[%s] runId=[%s] with metadata=[%+v]", job.ID, job.JobType, job.RunId, taskMetadata)
	}

	return tx.Commit().Error
}

// Runner Finds Jobs that havent been executed yet and schedule them
// (by creating db records)
func JobRunner(jobSpecStore JobSpecStore, db *gorm.DB) error {
	var jobs []JobRecord

	// find all job records that are created but not started running
	if err := db.
		Where("status = ?", Created).
		Where("start_time < ?", time.Now()).
		Find(&jobs).
		Error; err != nil {
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

func TaskRunner(jobSpecStore JobSpecStore, db *gorm.DB) error {
	var tasks []TaskRecord
	if err := db.
		Where("status = ?", Created).
		Find(&tasks).
		Error; err != nil {
		return err
	}

	syncChannel := make(chan TaskRecord)

	numTasks := len(tasks)

	for _, task := range tasks {
		rlog.Infof("Running task %d: %#v", task.ID, task)
		go RunTask(db, syncChannel, jobSpecStore, task)
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
