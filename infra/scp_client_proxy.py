#!/usr/bin/python3

from scp import SCPClient
from infra.ssh_client_proxy import SSHClientProxy

import logging

logger = logging.getLogger(__name__)
logger.setLevel(logging.debug)
ch = logging.StreamHandler(sys.stderr)
logger.addHandler(ch)


class SCPClientProxy(object):
    def __init__(self, ssh_client_proxy: SSHClientProxy):
        logger.debug("Creating new scp client")
        self.scp_client = SCPClient(
            ssh_client_proxy.get_ssh_client().get_transport(),
        )

    def put(self, *args, **kwargs):
        self.scp_client(*args, **kwargs)

    def __enter__(self):
        return self

    def __exit__(self):
        self.scp_client.close()

    def get_scp_client(self):
        return self.scp_client
