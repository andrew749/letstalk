#!/bin/sh

set -e
set -x

echo "Starting job runner"
go run jobs/job_runner/main.go
echo "Finished job runner"

echo "Starting task runner"
go run jobs/task_runner/main.go
echo "Finished task runner"

echo "Starting status runner"
go run jobs/job_status_updater/main.go
echo "Finished status runner"
