#!/bin/bash

docker run -it letstalk_webapp bash <<'EOF'
  SECRETS_PATH="secrets.json"
  # get options

  echo "Getting constants"
  export DB_PASS=$(jq -r '.db_pass' ${SECRETS_PATH})
  export DB_ADDR=$(jq -r '.db_addr' ${SECRETS_PATH})
  export DB_USER=$(jq -r '.db_user' ${SECRETS_PATH})

  if [[ -z "$DB_PASS" ]]; then
    echo "Database password not set."
    exit 1
  fi

  if [[ -z "$DB_ADDR" ]]; then
    echo "Database url not set."
    exit 1
  fi

  if [[ -z "$DB_USER" ]]; then
    echo "Database user not set."
    exit 1
  fi

  echo "Got constants"

  # run the repl
  gore <<EOF
    :import letstalk/server/utility
  EOF
EOF
