#!/bin/bash

PASSWORD=$(jq -r '.db_pass' server/secrets.json)
mysql -h letstalk.cobpajnfrzy1.us-east-1.rds.amazonaws.com -u hiveadmin -p$PASSWORD
