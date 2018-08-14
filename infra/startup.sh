#!/bin/sh
# script to run on startup
APP=/var/app/letstalk/infra
SCRIPT=deploy.sh

cd $APP && ./$SCRIPT
