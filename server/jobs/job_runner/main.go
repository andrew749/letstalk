package main

import (
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs"
	"letstalk/server/utility"

	"github.com/romana/rlog"
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		rlog.Errorf("Unable to get database: %+v", err)
		panic(err)
	}

	// create new job runner to run jobs
	err = jobmine.JobRunner(jobmine_jobs.Jobs, db)
	if err != nil {
		rlog.Errorf("Job runner ran into exception: %+v", err)
		panic(err)
	}
}
