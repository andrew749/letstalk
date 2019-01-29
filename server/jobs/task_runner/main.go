package main

import (
	"flag"
	"letstalk/server/constants"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs"
	sqs_notification_processor "letstalk/server/jobs/sqs_notification_processor/src"
	"letstalk/server/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/romana/rlog"
)

// NOTE: flags are conjunctive as in the conditions are ANDed together
var (
	jobTypeFilter = flag.String("jobTypeFilter", "", "Job Types to run")
	runIdFilter   = flag.String("runIdFilter", "", "Run Id to run")
)

func main() {
	utility.Bootstrap()
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

	rlog.Debugf("Queue processing")
	// process anything in the sqs queue
	if helper, ok := utility.QueueHelper.(utility.LocalQueueImpl); ok {
		helper.SubscribeListener(constants.NotificationQueueUrl, func(event *events.SQSEvent) error {
			if err := sqs_notification_processor.SendNotificationLambda(*event); err != nil {
				rlog.Critical(err)
				return err
			}
			return nil
		})
		go helper.QueueProcessor()

		// create new task runner
		rlog.Debug("Running tasks")
		err = jobmine.TaskRunner(jobmine_jobs.Jobs, db, jobTypes, runIds)
		if err != nil {
			rlog.Errorf("Task runner ran into exception: %+v", err)
			panic(err)
		}
		rlog.Debug("Done running tasks")

		helper.CloseQueue()
		rlog.Debugf("Finishing queue")
		helper.WaitForQueueDone()
		rlog.Debugf("Queue done processing")
	}
}
