#!/bin/bash

set -e

echo "Starting cluster"
docker-compose -f '../docker-compose.yml' -f '../docker-compose-prod.yml' up -d
