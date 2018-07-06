#!/bin/bash

DIRECTORY="/var/app/letstalk"
DB="letstalk.cobpajnfrzy1.us-east-1.rds.amazonaws.com"
USERNAME="hiveadmin"
OUT_DIR="/var"

PASSWORD=$(jq -r '.db_pass' $DIRECTORY/server/secrets.json)

mysql -h $DB  -u $USERNAME -p$PASSWORD -d > $OUT_DIR/$(date +%Y-%m-%dT%::z).db.sql
