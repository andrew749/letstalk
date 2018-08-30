#!/bin/bash

FOLDER=/var/app/letstalk
set -e
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
    sudo add-apt-repository ppa:certbot/certbot
    sudo apt-get update
    sudo apt-get install docker \
      docker-compose \
      jq \
      software-properties-common \
      python-certbot-nginx \
      apt-transport-https
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

setup_startup() {
  cp $FOLDER/infra/server /etc/init.d/server
  update-rc.d server defaults
}

setup_datadog() {
  sudo sh -c "echo 'deb https://apt.datadoghq.com/ stable 6' > /etc/apt/sources.list.d/datadog.list"
  sudo apt-key adv --recv-keys --keyserver hkp://keyserver.ubuntu.com:80 382E94DE

  # install agent
  sudo apt-get update
  sudo apt-get install datadog-agent

  # configure agent
  read -p "Datadog api key: " DATADOG_API_KEY
  sudo sh -c "sed 's/api_key:.*/api_key: $DATADOG_API_KEY/' /var/app/letstalk/infra/config/datadog.yaml > /etc/datadog-agent/datadog.yaml"
  systemctl start datadog-agent
  systemctl enable datadog-agent
}

# start of actual program
create_admin_group
create_admin_user
install_dependencies
generate_ssh
setup_docker
setup_datadog

echo "\033[92mAdding source code.\033[0m"
git clone git@github.com:andrew749/letstalk.git

echo "\033[91mRemember to manually add secrets to this server in the server root!!!\033[0m"
