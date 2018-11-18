#!/bin/bash
set -x

declare -a JOBS_TO_BUILD=(
  "task_runner"
  "job_runner"
  "job_status_updater"
  "manual_job_scheduler"
  "manual_task_scheduler"
)

JOB_DIR="jobs"
ENTRYPOINT="main.go"
ARTIFACT_PATH="build"

for job in "${JOBS_TO_BUILD[@]}";
do
  echo "Building $job"
  OUTPUT_PATH="$ARTIFACT_PATH/$job"
  GOOS=linux go build -o $OUTPUT_PATH "$JOB_DIR/$job/$ENTRYPOINT"
  echo "Built $job in $OUTPUT_PATH"
done
