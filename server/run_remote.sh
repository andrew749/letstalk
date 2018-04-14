#!/bin/bash

python3 run_local.py --db_addr="tcp(letstalk.cobpajnfrzy1.us-east-1.rds.amazonaws.com:3306)" --db_user="hiveadmin" --db_pass="$1"
