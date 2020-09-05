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
	jobTypeRaw = flag.String("jobType", "", "Job Type for this job")
	runId      = flag.String("runId", "", "Run id for this job")
	metadata   = flag.String("metadata", "", "JSON formatted metadata to pass to the job")
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

	var parsedMetadata jobmine.Metadata

	if err := json.Unmarshal([]byte(*metadata), &parsedMetadata); err != nil {
		rlog.Errorf("Failed to parse metadata.")
		panic(err)
	}

	// actually create the job.
	if _, err := jobmine.CreateJobRecord(db, *runId, jobType, parsedMetadata, nil); err != nil {
		rlog.Errorf("Failed to schedule reminder job.")
		panic(err)
	}
	rlog.Infof("Successfully scheduled reminder job.")
}
