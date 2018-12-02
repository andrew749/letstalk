#!/bin/bash

set -e
set -x

EXEC_COMMAND="source source_secrets.sh secrets.json && $@"
docker run --network="$DB_NET" -v "$(pwd):/go/src/letstalk/server" -it letstalk_env bash -c "$EXEC_COMMAND"
