#!/bin/bash

set -e
set -x

docker run --network="$DB_NET" -v "$(pwd):/go/src/letstalk/server" -it letstalk_env bash -c "source source_secrets.sh secrets.json && $@"
