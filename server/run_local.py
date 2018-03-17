#!/usr/bin/python3

import sys
from argparse import ArgumentParser
from subprocess import run
import subprocess
import os
import logging

#logging
logger=logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

# server configuration
GIN_MODE="release"
RLOG_LOG_LEVEL="DEBUG"
RLOG_TIME_FORMAT='2006-01-02T15:04:05'
DB_ADDR='tcp(127.0.0.1:3306)'
DB_USER='letstalk'
DB_PASS="uwletstalk"

# go env variables
GOPATH = os.environ['GOPATH']
GO_BINARY = "/usr/bin/go"
SERVER=f"{GOPATH}/src/letstalk/server"

def usage():
  print (
    'Usage : SECRETS_PATH=<path to secrets> ./run_local.py',
    file=sys.stderr,
  )
  sys.exit(1)

def get_args():
    parser = ArgumentParser(description="Run an instance of Hive")
    parser.add_argument(
        "--prod",
        action="store_true",
        help="whether to run the server in production mode",
    )
    return parser.parse_args()

def main():
    args = get_args()
    env = os.environ.update({
        "DB_ADDR": DB_ADDR,
        "DB_USER": DB_USER,
        "DB_PASS": DB_PASS,
        "RLOG_LOG_LEVEL": RLOG_LOG_LEVEL,
        "GOPATH": GOPATH,
    })

    # install dependencies
    if args.prod:
        logger.info("Running production server")
        logger.debug(env)
        process = subprocess.Popen(
            [f'{GO_BINARY}', 'run', 'core/main.go'],
            env=env,
        )
    else:
        logger.info("Running Debug Server")
        process = subprocess.Popen(
            [f"{GOPATH}/bin/gin", '--build', 'core', '--excludeDir', 'vendor'],
            env=env,
        )
    # wait for program to finish
    process.communicate()

if __name__ == "__main__":
    main()

