#!/usr/bin/python3

import sys
from argparse import ArgumentParser
from subprocess import run
from shutil import copytree
from utils.utils import template_file
import subprocess
import os
import logging

"""
TO BE RUN INSIDE DOCKER CONTAINER
"""

#logging
logger=logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

# server configuration
GIN_MODE="release"
RLOG_LOG_LEVEL="DEBUG"
RLOG_TIME_FORMAT='2006-01-02T15:04:05'

# default database parameters
DB_ADDR='tcp(localhost:3306)'
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
        help="Whether to run the server in production mode",
    )
    parser.add_argument(
        "--db_addr",
        default=DB_ADDR,
        help="The address of the db to connect to",
    )
    parser.add_argument(
        "--db_user",
        default=DB_USER,
        help="Username to connect to db with",
    )
    parser.add_argument(
        "--db_pass",
        default=DB_PASS,
        help="Password to connect to db with",
    )
    return parser.parse_args()

def main():
    args = get_args()
    env = os.environ.update({
        "DB_ADDR": args.db_addr,
        "DB_USER": args.db_user,
        "DB_PASS": args.db_pass,
        "RLOG_LOG_LEVEL": RLOG_LOG_LEVEL,
        "GOPATH": GOPATH,
    })

    # install dependencies
    if args.prod:
        logger.info("Running Production server")
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

