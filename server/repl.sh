#!/bin/bash

# get options
export DB_PASS=$(jq -r '.db_pass' ${SECRETS_PATH})
export DB_ADDR=$(jq -r '.db_addr' ${SECRETS_PATH})
export DB_USER=$(jq -r '.db_user' ${SECRETS_PATH})

gore

