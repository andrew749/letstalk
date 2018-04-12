#!/usr/bin/python3

from paramiko import SSHClient, RSAKey
import paramiko
from typing import List

import logging

logger = logging.getLogger(__name__)
logger.setLevel(logging.debug)
ch = logging.StreamHandler(sys.stderr)
logger.addHandler(ch)

class SSHClientProxy(object):
    def __init__(
        self,
        server_ip: str,
        username: str,
        private_key: RSAKey,
        port=22,
    ):
        self.client = SSHClient()
        self.client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        self.client.connect(server, port, username, pkey=private_key)

    def get_ssh_client(self) -> SSHClient:
        return self.client

    def run(self, command_list: List[str]):
        command = " ".join(command_list)
        logger.debug(command)
        return self.client.exec_command(command)

    def __enter__(self) -> SSHClientProxy:
        return self

    def __exit__(self):
        self.client.close()

    @staticmthod
    def load_private_key(private_key_path: str, password: str) -> RSAKey:
        return RSAKey.from_private_key_file(
            private_key_path,
            password=password,
        )
