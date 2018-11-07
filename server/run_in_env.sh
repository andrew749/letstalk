#!/bin/bash

docker run -v "$(pwd):/go/src/letstalk/server" -it letstalk_env bash -c "source source_secrets.sh secrets.json && $@"
