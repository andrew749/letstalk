#!/bin/bash

# Script to run the server pointing to a remote host
# Requires SECRETS_PATH env variable to be set
if [[ -z  "$SECRETS_PATH" ]]; then
    echo "Please set SECRETS_PATH"
    exit 1
fi

DB_PASS=$(jq -r '.db_pass' ${SECRETS_PATH})
if [[ -z "$DB_PASS" ]]; then
    echo "Database password not set."
    exit 1
fi

DB_ADDR=$(jq -r '.db_addr' ${SECRETS_PATH})
if [[ -z "$DB_ADDR" ]]; then
    echo "Database url not set."
    exit 1
fi

DB_USER=$(jq -r '.db_user' ${SECRETS_PATH})
if [[ -z "$DB_USER" ]]; then
    echo "Database user not set."
    exit 1
fi

ES_ADDR=$(jq -r '.elastic_search_url' ${SECRETS_PATH})
if [[ -z "ES_ADDR" ]]; then
    echo "Database user not set."
    exit 1
fi


if [[ -z $PROD ]]; then
    echo "REMOTE: running dev"
    python3 run_local.py
else
    echo "REMOTE: running prod"
    python3 run_local.py --db_addr="$DB_ADDR" --db_user="$DB_USER" \
      --db_pass="$DB_PASS" --es_addr="$ES_ADDR" --prod
fi
