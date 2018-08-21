#!/usr/bin/python3

import sys
from argparse import ArgumentParser
from subprocess import run
from shutil import copytree
import subprocess
import os
import logging

"""
TO BE RUN INSIDE DOCKER CONTAINER
"""

#logging
logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

# server configuration
GIN_MODE_DEBUG="debug"
GIN_MODE_PROD="release"
RLOG_LOG_LEVEL="DEBUG"
RLOG_TIME_FORMAT='2006-01-02T15:04:05'

# default database parameters
DB_ADDR='tcp(mysql:3306)'
DB_USER='letstalk'
DB_PASS="uwletstalk"

# default elasticsearch parameters
ES_ADDR='http://elasticsearch:9200'

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
    parser.add_argument(
        "--es_addr",
        default=ES_ADDR,
        help="The address of the elasticsearch cluster to connect to",
    )
    return parser.parse_args()

def main():
    args = get_args()
    os.environ.update({
        "DB_ADDR": args.db_addr,
        "DB_USER": args.db_user,
        "DB_PASS": args.db_pass,
        "ES_ADDR": args.es_addr,
        "RLOG_LOG_LEVEL": RLOG_LOG_LEVEL,
        "GOPATH": GOPATH,
        "GIN_MODE": GIN_MODE_PROD if args.prod else GIN_MODE_DEBUG,
    })

    # install dependencies
    if args.prod:
        logger.info("LOCAL: Running Production server")
        process = subprocess.Popen(
            [f'{GO_BINARY}', 'run', 'core/main.go'],
            env=os.environ,
        )
    else:
        logger.info("LOCAL: Running Debug Server")
        process = subprocess.Popen(
            [f"{GOPATH}/bin/gin", '--build', 'core', '--excludeDir', 'vendor'],
            env=os.environ,
        )
    # wait for program to finish
    process.communicate()

if __name__ == "__main__":
    main()
