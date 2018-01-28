#! /bin/sh

SERVER=$GOPATH/src/letstalk/server
go run $SERVER/core/main.go --db-user='letstalk' --db-pass='uwletstalk' --db-addr='tcp(127.0.0.1:3306)'

