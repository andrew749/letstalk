package main

import (
	"letstalk/server/jobmine_jobs/remind_onboard_job"
	"letstalk/server/utility"

	"github.com/romana/rlog"
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}
	rlog.Infof("Scheduling reminder job.")
	err = remind_onboard_job.CreateReminderJob(db)
	if err != nil {
		panic(err)
	}
	rlog.Infof("Successfully scheduled reminder job.")
}
