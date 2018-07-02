#!/bin/bash

set -e

echo "Building new version of application"
docker-compose -f '../docker-compose.yml' -f '../docker-compose-prod.yml' build webapp
