#!/bin/bash
set -e

docker exec -it nginx nginx -s reload
