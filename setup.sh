#!/bin/bash

# TO BE RUN ON SERVER RUNNING DOCKER

# setup latest git hooks
cp hooks/* .git/hooks
PROJECT_ROOT=$(pwd)
MOUNT_ROOT=/go/src

echo "Building docker image"
docker build -t hive:latest .

if [$? -ne 0]; then
    echo "Unable to build image."
    exit 1
fi

echo "Starting Container"
docker run -it -p 80:3000 -v $(pwd)/server:/go/src/letstalk/server hive:latest
