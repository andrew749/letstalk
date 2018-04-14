#!/bin/bash

# To be run from AWS cloud in an EC2 instance

# Group to create for server administration (including running the server)
ADMINGROUP="serveradmingrp"
ADMINUSER="serveradmin"

# directories for the app
HOME=/var/app/letstalk
SERVER=${HOME}/server

# create group and add the current user to the group
create_admin_group() {
    sudo groupadd $ADMINGROUP
}

create_admin_user() {
    # create user  without home directory that cant login
    sudo useradd -M $ADMINUSER
    sudo usermod -L $ADMINUSER
    sudo usermod -aG $ADMINGROUP $ADMINUSER
}

# install dependencies
install_dependencies() {
    sudo apt-get install docker
}

setup_docker() {
    # initialize the swarm
    docker swarm init
}


# start of actual program
create_admin_group
create_admin_user
install_dependencies
setup_docker

echo "Remember to manually add secrets to this server"
