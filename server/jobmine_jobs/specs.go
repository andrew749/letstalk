package jobmine_jobs

import (
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs/remind_meetup_job"
	"letstalk/server/jobmine_jobs/remind_onboard_job"
	"letstalk/server/jobmine_jobs/test_job"
	"letstalk/server/jobmine_jobs/welcome_back_email_job"
)

// IMPORTANT: JobType and the key need to be identical here.
// This defines the jobs that can get run and performs the mapping between values
// stored in the database to the actual code that should get executed.
var Jobs = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		test_job.TestJob:                              test_job.TestJobSpec,
		remind_onboard_job.REMIND_ONBOARD_JOB:         remind_onboard_job.ReminderJobSpec,
		remind_meetup_job.REMIND_MEETUP_JOB:           remind_meetup_job.ReminderJobSpec,
		welcome_back_email_job.WELCOME_BACK_EMAIL_JOB: welcome_back_email_job.ReminderJobSpec,
	},
}
