package main

import (
	"encoding/json"
	"letstalk/server/jobmine_jobs/test_job"
	"letstalk/server/utility"

	"github.com/namsral/flag"

	"github.com/romana/rlog"
)

var (
	runId       = flag.String("runId", "", "Unique identifier for this job run.")
	rawMetadata = flag.String("metadata", "", "Json metadata that is defined for the job")
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		rlog.Errorf("Unable to get db: %+v", err)
		panic(err)
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(*rawMetadata), &metadata)
	if err != nil {
		rlog.Errorf("Bad input metadata")
		panic(err)
	}

	err = test_job.CreateTestJob(db, *runId, metadata)
	if err != nil {
		panic(err)
	}
}
