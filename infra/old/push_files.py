#!/usr/bin/python3

from argparse import ArgumentParser
from ssh_client_proxy import SSHClientProxy
from scp_client_proxy import SCPClientProxy
from pusher import Pusher, PushCommand, ServerCommand
import boto3

import os
import logging
import getpass
import sys

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler(sys.stderr)
logger.addHandler(ch)

def get_args():
    parser = ArgumentParser(
        description="Send files to ec2 instances that we own.",
    )

    parser.add_argument(
        "files",
        nargs="+",
        help="Files to push"
    )

    parser.add_argument(
        "--destination",
        required=True,
        help="Destination folder to send files to."
    )

    parser.add_argument(
        "--host",
        default=None,
        help="Host to send files to."
    )

    # arguments to authenticate with server
    parser.add_argument(
        "--username",
        default=os.environ['USER'],
        help="The username to use to logon to host with."
    )

    parser.add_argument(
        "--private_key_path",
        default="~/.ssh/id_rsa",
        help="Private key to use to authenticate to host.",
    )

    return parser.parse_args()


def main():
    args = get_args()
    # get all hosts
    session = boto3.Session(profile_name='hive')
    client = session.resource('ec2', region_name="us-east-1")

    private_key_path = os.path.abspath(os.path.expanduser(args.private_key_path))

    # load private key from local filesystem
    password = getpass.getpass("Private Key Password to push with")
    private_key = SSHClientProxy.load_private_key(private_key_path, password)

    pusher = Pusher(
        username=args.username,
        private_key=private_key,
        hosts = [instance.public_ip_address for instance in client.instances.all()] if not args.host else args.host,
    )

    pusher.add_command(
        PushCommand(
            target_file_paths=args.files,
            destination=args.destination,
            post_copy_payload=ServerCommand(
                command=["ls" "-l"],
            ),
        )
    )

    logger.info("Pushing to all instances")
    pusher.push_all_hosts()
    logger.info("Successfully pushed to all instances")

if __name__ == "__main__":
    main()

