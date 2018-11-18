#!/bin/sh

set -e
set -x

if [ ! -f build/job_runner ] || [ ! -f build/task_runner ] || [ ! -f build/job_status_updater ];
then
  echo "Required binaries not found. Building from source ..."
  ./rebuild_job_scheduling_toolchain.sh
fi

echo "Starting job runner"
./build/job_runner
echo "Finished job runner"

echo "Starting task runner"
./build/task_runner
echo "Finished task runner"

echo "Starting status runner"
./build/job_status_updater
echo "Finished status runner"
