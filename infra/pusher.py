#!/usr/bin/python3

from ssh_client_proxy import SSHClientProxy, RSAKey
from scp_client_proxy import SCPClientProxy

from typing import List, Tuple

import logging
import sys

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler(sys.stderr)
logger.addHandler(ch)


class RemoteCommand(object):
    def get_description(self):
        raise NotImplementedError()

    def execute(self, ssh_client: SSHClientProxy):
        raise NotImplementedError()


class ServerCommand(RemoteCommand):
    """
    Abstraction representing arbitrary commands that can get executed on servers
    """
    def __init__(self, command: List[str]) -> None:
        self.command = command

    def get_description(self):
        return "Command: {}".format(" ".join(self.command))

    def execute(self, ssh_client: SSHClientProxy) -> Tuple[int, int, int] :
        return ssh_client.run(self.command)


class PushCommand(RemoteCommand):
    """
    Class to encapsulate what actually gets copied to a server.
    """
    def __init__(
        self,
        target_file_paths: List[str],
        destination: str,
        post_copy_payload: ServerCommand,
    ) -> None:
        self.target_file_paths = target_file_paths
        self.destination = destination
        self.post_copy_payload = post_copy_payload

    def get_description(self):
        return "Pushing {} to {}".format(self.target_file_paths, self.destination)

    def execute(self, ssh_client: SSHClientProxy):
        with SCPClientProxy(ssh_client) as scp_client:
            logging.debug("Pushing file")
            scp_client.put(
                self.target_file_paths,
                self.destination,
                recursive=True,
            )
            self.post_copy_payload.execute(ssh_client)


class Pusher(object):
    """
    Class to facilitate pushing multiple payloads to servers.
    """

    def __init__(
        self,
        username: str,
        private_key: RSAKey,
        hosts: List[str] = [],
        commands: List[RemoteCommand] = [],
    ) -> None:
        self.hosts = hosts
        self.commands = commands
        self.username = username
        self.private_key = private_key

    def add_host(self, host):
        self.hosts.append(host)

    def add_command(self, command: RemoteCommand):
        """
        Add to the commands that should be executed on each server.
        """
        self.commands.append(command)

    def push_host(self, host: str):
        # create ssh client for this host
        logger.debug("Creating ssh client for host %s", host)
        ssh_client = SSHClientProxy(
            server_ip=host,
            username=self.username,
            private_key=self.private_key,
        )
        # execute all push commands
        for command in self.commands:
            logger.debug(command.get_description())
            command.execute(ssh_client)

    def push_all_hosts(self):
        for host in self.hosts:
            logger.debug("Pushing host: %s", host)
            self.push_host(host)

