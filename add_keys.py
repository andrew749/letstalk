#!/usr/bin/python3

from argparse import ArgumentParser
import boto3
from paramiko import SSHClient
from scp import SCPClient

import paramiko

import os
import logging
import getpass
import sys

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler(sys.stdout)
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

def load_private_key(private_key_path: str, password: str):
    return paramiko.RSAKey.from_private_key_file(
        private_key_path,
        password=password,
    )

def createSSHClient(server: str, username: str, private_key, port=22):
    client = SSHClient()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    client.connect(server, port, username, pkey=private_key)
    return client

def createSCPClient(ssh: SSHClient):
    return SCPClient(ssh.get_transport())

def create_user_if_not_exists(username: str, ssh_client: SSHClient):
    # need password so we can do sudo
    password = getpass.getpass()
    (stdin, stdout, stderr) = execute_server_cmd(
        "sudo useradd -m {}".format(username),
        ssh_client,
    )
    stdin.write(password)
    # add directory for ssh
    (stdin, stdout, stderr) = execute_server_cmd(
        "sudo mkdir -p {}".format(os.path.join("/home", username, ".ssh")),
        ssh_client,
    )
    stdin.write(password)


def execute_server_cmd(command: str, ssh_client: SSHClient):
    logger.debug(command)
    return ssh_client.exec_command(command)

def provision_server(username_to_push, admin_username, server, private_key, public_key):
    ssh_client = createSSHClient(server, admin_username, private_key)
    with createSCPClient(ssh_client) as scp_client:
        # add user if it doesnt exist
        create_user_if_not_exists(username_to_push, ssh_client)

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
    password = getpass.getpass()
    private_key = load_private_key(private_key_path, password)

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

