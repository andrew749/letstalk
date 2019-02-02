package main

import (
	"encoding/json"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs"
	"letstalk/server/utility"

	"github.com/namsral/flag"

	"github.com/romana/rlog"
)

var (
	runId        = flag.String("runId", "", "Run id for this task")
	jobTypeRaw   = flag.String("jobType", "", "Job Type for this job")
	jobMetadata  = flag.String("jobMetadata", "", "JSON formatted metadata to pass to the job")
	taskMetadata = flag.String("taskMetadata", "", "JSON formatted metadata to pass to the task")
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}

	jobType := jobmine.JobType(*jobTypeRaw)
	// check that the job is valid
	jobSpecStore := jobmine_jobs.Jobs
	_, err = jobSpecStore.GetJobSpecForJobType(jobType)
	if err != nil {
		rlog.Errorf("Unable to find job %s", *jobTypeRaw)
		panic(err)
	}

	rlog.Infof("Scheduling job.")

	var (
		parsedJobMetadata  jobmine.Metadata
		parsedTaskMetadata jobmine.Metadata
	)

	if err := json.Unmarshal([]byte(*jobMetadata), &parsedJobMetadata); err != nil {
		rlog.Errorf("Failed to parse metadata.")
		panic(err)
	}

	if err := json.Unmarshal([]byte(*taskMetadata), &parsedTaskMetadata); err != nil {
		rlog.Errorf("Failed to parse metadata.")
		panic(err)
	}
	tx := db.Begin()

	jobRecord, err := jobmine.CreateJobRecord(tx, *runId, jobType, parsedJobMetadata, nil)
	// actually create the job.
	if err != nil {
		tx.Rollback()
		rlog.Errorf("Failed to create job")
		panic(err)
	}
	rlog.Infof("Successfully created job.")

	taskRecord, err := jobmine.CreateTaskRecord(tx, jobRecord.ID, *runId, jobType, parsedTaskMetadata)
	if err := db.Create(&taskRecord).Error; err != nil {
		tx.Rollback()
		rlog.Errorf("Failed to create task")
		panic(err)
	}

	rlog.Infof("Successfully created task.")
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
	rlog.Info("Succesfully commited transaction.")
}
