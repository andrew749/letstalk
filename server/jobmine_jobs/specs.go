package jobmine_jobs

import (
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs/generic_notification_job"
	"letstalk/server/jobmine_jobs/match_round_commit_job"
	"letstalk/server/jobmine_jobs/notification_status_checker"
	"letstalk/server/jobmine_jobs/remind_meetup_job"
	"letstalk/server/jobmine_jobs/remind_onboard_job"
	"letstalk/server/jobmine_jobs/seed_mentorships_job"
	"letstalk/server/jobmine_jobs/test_job"
	"letstalk/server/jobmine_jobs/welcome_back_email_job"
)

// IMPORTANT: JobType and the key need to be identical here.
// This defines the jobs that can get run and performs the mapping between values
// stored in the database to the actual code that should get executed.
var Jobs = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		test_job.TestJob:                                            test_job.TestJobSpec,
		remind_onboard_job.REMIND_ONBOARD_JOB:                       remind_onboard_job.ReminderJobSpec,
		remind_meetup_job.REMIND_MEETUP_JOB:                         remind_meetup_job.ReminderJobSpec,
		seed_mentorships_job.SEED_MENTORSHIPS_JOB:                   seed_mentorships_job.SeedJobSpec,
		welcome_back_email_job.WELCOME_BACK_EMAIL_JOB:               welcome_back_email_job.ReminderJobSpec,
		notification_status_checker.NOTIFICATION_STATUS_CHECKER_JOB: notification_status_checker.NotificationStatusChecker,
		generic_notification_job.GENERIC_NOTIFICATION_JOB:           generic_notification_job.GenericNotificationJob,
		match_round_commit_job.MATCH_ROUND_COMMIT_JOB:               match_round_commit_job.CommitJobSpec,
	},
}
