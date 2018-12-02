#!/bin/bash

set -e
set -x

EXEC_COMMAND="source source_secrets.sh secrets.json && $@"
if [ -n "$PROD" ]; then
    VOL_DIR="/var/app/letstalk/server"
else
    VOL_DIR=$(pwd)
fi
docker run --network="$DB_NET" -v "$VOL_DIR:/go/src/letstalk/server" -it letstalk_env bash -c "$EXEC_COMMAND"
