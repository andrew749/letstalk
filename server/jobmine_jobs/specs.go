package jobmine_jobs

import (
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs/test_job"
)

// IMPORTANT: JobType and the key need to be identical here.
// This defines the jobs that can get run and performs the mapping between values
// stored in the database to the actual code that should get executed.
var Jobs = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		test_job.TestJob: test_job.TestJobSpec,
	},
}
