#!/bin/bash

# Used for development environment.

# setup latest git hooks
cp hooks/* .git/hooks

docker-compose -f docker-compose-debug.yml up
