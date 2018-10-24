package main

import (
	"letstalk/server/jobmine"
	"letstalk/server/utility"

	"github.com/romana/rlog"
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		rlog.Errorf("Unable to get database: %+v", err)
		panic(err)
	}

	err = jobmine.JobStateWatcher(db)
	if err != nil {
		rlog.Errorf("JobStateWatcher  ran into exception: %+v", err)
		panic(err)
	}
}
