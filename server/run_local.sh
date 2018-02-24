#! /bin/sh

export GIN_MODE=release
export RLOG_LOG_LEVEL=DEBUG
export RLOG_TIME_FORMAT='2006-01-02T15:04:05'
#RLOG_CALLER_INFO=true
#RLOG_LOG_NOTIME
#RLOG_LOG_FILE
#RLOG_LOG_STREAM

usage ()
{
  echo 'Usage : SECRETS_PATH=<path to secrets> ./run_local.sh'
  exit
}

if [[ -z $SECRETS_PATH ]];
then
  usage
fi

SERVER=$GOPATH/src/letstalk/server
DB_USER='letstalk' DB_PASS="uwletstalk" DB_ADDR='tcp(127.0.0.1:3306)' \
  $GOPATH/bin/gin --build="core" run core/main.go
