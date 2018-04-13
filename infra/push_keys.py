#!/usr/bin/python3

from argparse import ArgumentParser
from ssh_client_proxy import SSHClientProxy
from scp_client_proxy import SCPClientProxy
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
        description="Send public key certs to ec2 instances that we own.",
    )

    parser.add_argument(
        "--public_key_path",
        required=True,
        help="The public key path to send to the server",
    )

    parser.add_argument(
        "--username_to_push",
        required=True,
        help="The username to use to add this key for. If no home directory exists, create a new one."
    )

    # arguments to authenticate with server
    parser.add_argument(
        "--admin_username",
        required=True,
        help="The username to use to authenticate with the server."
    )

    parser.add_argument(
        "--private_key_path",
        required=True,
        help="Private key to use to authenticate to servers.",
    )

    return parser.parse_args()

def create_user_if_not_exists(username: str, ssh_client: SSHClientProxy):
    # need password so we can do sudo
    password = getpass.getpass()
    (stdin, stdout, stderr) = ssh_client.run(["sudo","useradd","-m", username])
    stdin.write(password)
    # add directory for ssh
    (stdin, stdout, stderr) = ssh_client.run(
        ["sudo", "mkdir", "-p", os.path.join("/home", username, ".ssh")]
    )
    stdin.write(password)

def provision_server(username_to_push, admin_username, server, private_key, public_key):
    with SSHClientProxy(server, admin_username, private_key) as ssh_client_proxy:
        with SCPClientProxy(ssh_client_proxy) as scp_client:
            # add user if it doesnt exist
            create_user_if_not_exists(username_to_push, ssh_client_proxy)

            # add user directory
            scp_client.put(
                public_key,
                remote_path=os.path.join("/home", username_to_push, ".ssh", 'authorized_keys'),
            )

def main():
    logger.debug("Adding keys...")
    args = get_args()
    # get all hosts
    session = boto3.Session(profile_name='hive')
    client = session.resource('ec2', region_name="us-east-1")

    private_key_path = os.path.abspath(os.path.expanduser(args.private_key_path))
    public_key_path = os.path.abspath(os.path.expanduser(args.public_key_path))
    # load private key from local filesystem
    password = getpass.getpass("Private Key Password to push with")
    private_key = SSHClientProxy.load_private_key(private_key_path, password)

    # for each instance
    for instance in client.instances.all():
        logger.info("Pushing to instance %s", instance.public_ip_address)
        # ssh into the instance and provision
        provision_server(
            username_to_push=args.username_to_push,
            admin_username=args.admin_username,
            server=instance.public_ip_address,
            private_key=private_key,
            public_key=public_key_path,
        )
    logger.info("Successfully pushed to all instances")

if __name__ == "__main__":
    main()

