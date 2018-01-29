#! /bin/sh

export GIN_MODE=release
export RLOG_LOG_LEVEL=DEBUG
export RLOG_TIME_FORMAT='2006-01-02T15:04:05'
#RLOG_CALLER_INFO=true
#RLOG_LOG_NOTIME
#RLOG_LOG_FILE
#RLOG_LOG_STREAM

SERVER=$GOPATH/src/letstalk/server
go run $SERVER/core/main.go --db-user='letstalk' --db-pass='uwletstalk' --db-addr='tcp(127.0.0.1:3306)' "$@"

