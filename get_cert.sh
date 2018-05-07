#!/bin/bash
sudo apt-get install software-properties-common
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt-get install python-certbot-nginx
sudo certbot --nginx certonly

certbot certonly -d hiveapp.org -d www.hiveapp.org -d api.hiveapp.org -d beta.hiveapp.org
