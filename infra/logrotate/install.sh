#!/bin/bash

case $1 in
    --install)
    apt-get install logrotate
    ;;
    esac

# copy files to dest
cp *.logrotate /etc/logrotate.d/

# init config
logrotate -fv /etc/logrotate.d/*.logrotate
