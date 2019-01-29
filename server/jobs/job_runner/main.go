package main

import (
	"flag"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs"
	"letstalk/server/utility"

	"github.com/romana/rlog"
)

// NOTE: flags are conjunctive as in the conditions are ANDed together
var (
	jobTypeFilter = flag.String("jobTypeFilter", "", "Job Types to run")
	runIdFilter   = flag.String("runIdFilter", "", "Run Id to run")
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		rlog.Errorf("Unable to get database: %+v", err)
		panic(err)
	}

	var jobTypes []jobmine.JobType
	var runIds []string

	if jobTypeFilter != "" {
		rawTypes = utility.ParseListCommandLineFlags(jobTypeFilter)
		jobTypes = make([]jobmine.JobType, 0)
		for _, t := range rawTypes {
			jobTypes = append(jobTypes, jobmine.JobType(t))
		}
	}

	if runIdFilter != "" {
		runIds = utility.ParseListCommandLineFlags(runIdFilter)
	}

	// create new job runner to run jobs
	err = jobmine.JobRunner(jobmine_jobs.Jobs, db, jobTypes, runIds)
	if err != nil {
		rlog.Errorf("Job runner ran into exception: %+v", err)
		panic(err)
	}
}
