#!/bin/bash

# note only one person should be able to use this at a time since it uses static port allocation
# and the firewall needs a port open

# Also note that you need
# GatewayPorts yes
# set in the /etc/ssh/sshd_config

# $port must be in the range 10100-10200

read -p "Enter username to use for ssh: " username
read -p "Enter local port to forward(on this computer): " port
read -p "Select a port to use on the remote server[10100,10200]: " portRemote
ssh -g -R $portRemote:localhost:$port $username@hiveapp.org

echo "Now you can access your port $port at hiveapp.org:$portRemote! MAKE SURE TO KEEP THIS CONNECTION OPEN"
