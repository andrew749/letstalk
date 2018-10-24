# Jobmine

Jobmine is a job scheduling system that allows us to define new jobs to get run
asynchronously from core application logic.

Jobmine is backed by a db to keep track of the state of a job. Using this,
one can idempotently execute runs of jobs.


## Architecture
Jobs each represent a logical operation that should be performed en masse.
A task is an execution of a job with specific metadata that is customized for the execution.
Jobs also contain metadata that is given to all tasks which are instantiation of the job.

Specs define the code to get run and Records contain configuration and metadata specific
to an instantiation of a job.

A good example is how one would schedule a notification campaign.
A JobSpec with TaskSpec would define the specific code to get run for the notification. The campaign might have shared metadata that all tasks want to access (i.e. the job metadata). Each job is composed of tasks that must run (i.e. trying to send a notification to a specific user).

## IMPORTANT
Currently this job running system is designed under the assumption of a single
worker job. If this is not honored, then it is possible for things to get run more than once.

There are a few things that need to change to support more worker, notably:
- make the job watchers only take some jobs, not all the jobs
- once a job watcher gets a job, dont allow other nodes to get this job (maybe use zk for orchestration)
- symmetric issues for tasks

This will be fixed in a future update, but needs some sort of locking
to actually solve and therefore will be punted. Sorry.
