#!/bin/bash

# To be run from AWS cloud

HOME=/var/app/letstalk
SERVER=${HOME}/server

# initialize the swarm
docker swarm init

# build the containers and connect them
docker-compose up --build

