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
DB_ADDR='tcp(localhost:3306)'
DB_USER='letstalk'
DB_PASS="uwletstalk"

# go env variables
GOPATH = os.environ['GOPATH']
GO_BINARY = "/usr/bin/go"
SERVER=f"{GOPATH}/src/letstalk/server"
NGINX_CONFIG_FILE="hiveapp.nginx.conf"
NGINX_CONFIG_PATH=""
NGINX_INSTALL_PATH="/etc/nginx/sites-available/"

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
    return parser.parse_args()

def provision_nginx():
    """
    Install nginx ssl certs and restart the server
    """
    template_file(
        file_path=os.path.join(NGINX_CONFIG_PATH, NGINX_CONFIG_FILE),
        out_path=os.path.join(NGINX_INSTALL_PATH, NGINX_CONFIG_FILE),
        fill_in_dict={
        },
    )

    # copy ssl certs to folder
    os.makedirs("/etc/nginx/ssl/hiveapp")
    copytree("dev_certs/", "/etc/nginx/ssl/hiveapp/")

    # restart nginx
    run(["service" "nginx", "restart"])

def provision(is_prod=False):
    provision_nginx()

def main():
    args = get_args()
    env = os.environ.update({
        "DB_ADDR": DB_ADDR,
        "DB_USER": DB_USER,
        "DB_PASS": DB_PASS,
        "RLOG_LOG_LEVEL": RLOG_LOG_LEVEL,
        "GOPATH": GOPATH,
    })

    provision()
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

