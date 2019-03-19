#!/bin/bash

FOLDER=/var/app/letstalk
set -e
# To be run from AWS cloud in an EC2 instance

# directories for the app
APPDIR=/var/app
APP=${APPDIR}/letstalk
SERVER=${APP}/server
SECRETS_PATH=${SERVER}/secrets.json
DATADOG_CONF=/etc/datadog-agent/conf.d
DATADOG_DOCKER_CONF=$DATADOG_CONF/docker.d

# install dependencies
install_dependencies() {
    sudo add-apt-repository ppa:certbot/certbot
    sudo apt-get update
    sudo apt-get install docker \
      docker-compose \
      mysql-client \
      jq \
      software-properties-common \
      python-certbot-nginx \
      apt-transport-https \
      python3 \
      nodejs \
      npm \
      virtualenv \
      ruby

    sudo gem update
    sudo gem install mustache
}

setup_docker() {
    # initialize the swarm
    docker swarm init
}

generate_ssh() {
  sudo ssh-keygen
  echo "BEGIN PUBLIC KEY"
  cat /home/$USER/.ssh/id_rsa.pub
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
  sudo systemctl start datadog-agent
  sudo systemctl enable datadog-agent

  # install agent checks
  sudo cp $APP/infra/monitoring/docker_daemon.yaml $DATADOG_DOCKER_CONF

  #setup permissions
  sudo usermod -a -G docker dd-agent

  # restart agent
  sudo systemctl restart datadog-agent
}

# install logging service
install_logging() {
  # put the systemd service in the appropriate folder
  sudo bash "cat $SECRETS_PATH | mustache - $APP/infra/healthcheck/nginx_tailer.service > /lib/systemd/system/nginx_tailer.service"

  # install pip and dependencies
  pushd $APP/infra/healthcheck
    virtualenv -p /usr/bin/python3 .
    source ./bin/activate
    pip install -r requirements.txt
    deactivate
  popd
  # enable to service to start
  sudo systemctl enable nginx_tailer.service
}

install_dependencies

echo "Creating server administration user"
create_admin_group
create_admin_user
echo "DONE: Creating server administration user"

echo "Generating ssh keys"
generate_ssh

# wait for user input
read -n 1 -p "Add this ssh key to the github repo."
sudo mkdir -p $APPDIR
sudo chown $USER $APPDIR
cd $APPDIR
echo "\033[92mAdding source code.\033[0m"
git clone git@github.com:andrew749/letstalk.git


sudo groupadd docker
sudo usermod -aG docker $USER
# add the admin
sudo systemctl start docker

setup_docker

read -n 1 -p "\033[91mManually add secrets.json to this server in the server root $APP!!!\033[0m"
setup_datadog
install_logging

echo "Server is provisioned."

