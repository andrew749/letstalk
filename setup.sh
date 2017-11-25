#!/bin/sh

cp hooks/* .git/hooks
PROJECT_ROOT=$(pwd)
MOUNT_ROOT=/home/app

# setup database
if [ "$(uname)" == "Darwin" ]; then
    brew install postgresl redis
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    # Do something under GNU/Linux platform
    sudo apt-get install postgresql
fi

# build docker image
docker build .

# run the container and mount the project
docker run -it -p 8080:8080 -v $PROJECT_ROOT:$MOUNT_ROOT
