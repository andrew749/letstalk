#!/bin/bash
set -e

usage() {
echo << EOF
  Installs LetsEncrypt software and tries to get key for application's domain
EOF
  exit
}

sudo apt-get install software-properties-common
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt-get install python-certbot-nginx
sudo certbot --nginx certonly

certbot certonly -d hiveapp.org -d www.hiveapp.org -d api.hiveapp.org -d beta.hiveapp.org
