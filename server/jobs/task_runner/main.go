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
	// create new
	err = jobmine.TaskRunner(jobmine_jobs.Jobs, db)
	if err != nil {
		rlog.Errorf("Task runner ran into exception: %+v", err)
		panic(err)
	}
}
