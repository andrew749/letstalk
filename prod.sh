#!/bin/bash

echo "Starting cluster"
docker-compose up -f docker-compose.yml --build -d
