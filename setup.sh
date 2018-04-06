#!/bin/bash

# TO BE RUN ON SERVER RUNNING DOCKER

# read command line arguments
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -d|--debug)
    DEBUG_MODE=true
    shift # past argument
    ;;
esac
done

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

# DEBUG
debug() {
    docker run -it -p 80:3000 -v $(pwd)/server:/go/src/letstalk/server hive:latest
}

# PRODUCTION
production() {
    docker run --net="bridge" --add-host=localhost:`ip route show | grep docker0 | awk '{print \$9}'` -d -p 80:3000 hive:latest
}

# run the specific container
if [[$DEBUG]]; then
    debug
else
    production
fi
