#!/usr/bin/python3

from scp import SCPClient
from ssh_client_proxy import SSHClientProxy

import sys

class SCPClientProxy(SCPClient):
    def __init__(self, ssh_client_proxy: SSHClientProxy, *args, **kwargs) -> None:
        super().__init__(
            ssh_client_proxy.get_ssh_client().get_transport(),
            *args,
            **kwargs,
        )
