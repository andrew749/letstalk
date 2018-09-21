#!/bin/bash

# note only one person should be able to use this at a time since it uses static port allocation
# and the firewall needs a port open

read -p "Enter username: " username
read -p "Enter port to forward to: " port
ssh -g -R 10123:localhost:$port $username@hiveapp.org

echo "Now you can access your port $port at hiveapp.org:10123! MAKE SURE TO KEEP THIS CONNECTION OPEN"
