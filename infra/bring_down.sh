#!/bin/bash

set -e

usage(){
  echo "Usage $0 -p[roduction]"
}

production() {
  docker-compose -f docker-compose.yml -f docker-compose-prod.yml down
}

debug() {
  docker-compose -f docker-compose.yml -f docker-compose-debug.yml down
}

# run debug by default
if [[ -z $1 ]]; then
  debug
  exit
fi

while [ "$1" != "" ]; do
    case $1 in
        -p | --production )     production
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done
