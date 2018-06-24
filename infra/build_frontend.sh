#!/bin/bash

# Build files for the static site frontend.

# Get the apprpriate nodejs and gulp
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
apt-get install npm
npm install -g gulp

# setup nodejs environment
ln -s /usr/bin/nodejs /usr/bin/node

# install dependencies
npm --prefix ../landing install

# build
gulp --cwd ../landing
