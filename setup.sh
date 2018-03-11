#!/bin/sh

# setup latest git hooks
cp hooks/* .git/hooks
PROJECT_ROOT=$(pwd)
MOUNT_ROOT=/go/src

# setup database
if [ "$(uname)" == "Darwin" ]; then
    echo "Detected Mac"
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    echo "Detected Linux"
fi

echo "Building docker image"
docker build -t hive:latest .

echo "Starting Container"
# run the container and mount the project
docker run -it -p 3000:80 -v $(pwd)/server:/go/src/letstalk/server hive:latest
