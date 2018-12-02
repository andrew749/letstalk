#!/bin/bash

set -e
set -x

# Flags:
# - PROD: whether to use the cwd or the prod working directory
# - TTY: whether to allow tty input to docker

EXEC_COMMAND="source source_secrets.sh secrets.json && $@"
if [ -n "$PROD" ]; then
    VOL_DIR="/var/app/letstalk/server"
else
    VOL_DIR=$(pwd)
fi

TTY_ARGS=""
if [ -n "$TTY" ]; then
    TTY_ARGS="-it"
fi

docker run $TTY_ARGS --network="$DB_NET" -v "$VOL_DIR:/go/src/letstalk/server" letstalk_env bash -c "$EXEC_COMMAND"
