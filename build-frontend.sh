#!/bin/bash
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
apt-get install npm
npm install -g gulp

ln -s /usr/bin/nodejs /usr/bin/node
npm --prefix landing install
gulp --cwd landing
