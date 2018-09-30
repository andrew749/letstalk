#!/bin/bash

set -e
# Build files for the static site frontend.

FRONTEND_DIRECTORY="../landing"

# Get the apprpriate nodejs and gulp if not installed
NODE_INSTALLED=$(apt-cache policy nodejs | grep  Installed: | sed s/Installed://)
if [[ -z $NODE_INSTALLED ]]; then
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
apt-get install npm
npm install -g gulp

# setup nodejs environment
ln -s /usr/bin/nodejs /usr/bin/node
fi

# install dependencies
npm --prefix $FRONTEND_DIRECTORY install

# build
gulp --cwd $FRONTEND_DIRECTORY
