#!/bin/bash

# Used for development environment.

# setup latest git hooks
cp hooks/* .git/hooks

# if the environment image was already built, then rebuilt web assets
if [ $(docker images -q letstalk_webapp) -a $(docker images -q letstalk_env) ];
then
  # essentially mount our local directory web build directory into this container so we get the
  # build artifacts from docker
  docker run -it -v "$(pwd)"/server/web/dist:/go/src/letstalk/server/web/dist letstalk_env yarn --cwd web build
fi

docker-compose -f 'docker-compose.yml' -f 'docker-compose-debug.yml' up
