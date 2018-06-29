#!/bin/bash

# To be run from AWS cloud in an EC2 instance

# Group to create for server administration (including running the server)
ADMINGROUP="server_grp"
ADMINUSER="server"

# directories for the app
HOME=/var/app/letstalk
SERVER=${HOME}/server

# create group and add the current user to the group
create_admin_group() {
    sudo groupadd $ADMINGROUP
}

create_admin_user() {
    sudo useradd $ADMINUSER
    sudo usermod -aG $ADMINGROUP $ADMINUSER
}

# install dependencies
install_dependencies() {
    sudo apt-get install docker docker-compose jq
}

setup_docker() {
    # initialize the swarm
    docker swarm init
}

generate_ssh() {
  su $ADMINUSER
  ssh-keygen
  echo "BEGIN PUBLIC KEY"
  cat ~/.ssh/id_rsa.pub
  echo "END PUBLIC KEY"
}


# start of actual program
create_admin_group
create_admin_user
install_dependencies
generate_ssh
setup_docker

echo "\033[92mAdding source code.\033[0m"
git clone git@github.com:andrew749/letstalk.git

echo "\033[91mRemember to manually add secrets to this server\033[0m"