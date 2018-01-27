#! /bin/sh

SERVER=$GOPATH/src/letstalk/server

mysql -u letstalk -puwletstalk letstalk < $SERVER/db/letstalk.sql
cd $SERVER
modelq -db="modelq@/letstalk" -pkg=data -driver=mysql -schema=letstalk -p=8

