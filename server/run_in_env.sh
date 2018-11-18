#!/bin/bash

docker run --network="letstalk_db_net" -v "$(pwd):/go/src/letstalk/server" -it letstalk_env bash -c "source source_secrets.sh secrets.json && $@"
