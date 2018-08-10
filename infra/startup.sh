#!/bin/sh
# script to run on startup
FOLDER=/var/app/letstalk/infra
SCRIPT=deploy.sh

cd $FOLDER && ./$SCRIPT
