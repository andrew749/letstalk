#!/bin/bash

# build new nginx instance
docker build -t latest_load_balancer .

# run the load balancer
docker run --name load_balancer -d latest_load_balancer
