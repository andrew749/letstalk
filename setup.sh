#!/bin/sh

# setup latest git hooks
cp hooks/* .git/hooks
PROJECT_ROOT=$(pwd)
MOUNT_ROOT=/go/src

# setup database
if [ "$(uname)" == "Darwin" ]; then
    brew install mysql
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    # Do something under GNU/Linux platform
    sudo apt-get install mysql-server
fi

# build docker image
docker build -t hive:latest .

# run the container and mount the project
docker run -it -p 3000:80 -v server:/go/src/letstalk hive:latest
