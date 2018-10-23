# Hive (tentative)
## A mentorship platform
[![CircleCI](https://circleci.com/gh/andrew749/letstalk.svg?style=svg&circle-token=188ccb7b28649151618bf95dd0259cd67a5a1b9f)](https://circleci.com/gh/andrew749/letstalk)

## Basic project structure

`letstalk/`
This is the react native app. Code for iOS and Android clients lives here.

`server/`
This is the main messenger service. All backend code lives here.

`infra/`
Scripts to help with administration of servers


## Prerequisites
Install the following packages.
```
docker
docker-compose
go-dep
```

### How to see the dev database
Starting the server spins up its own mysql instance inside a docker container. If you see an error at startup, make sure to kill any already running instances of mysql.
```
mysql -h 127.0.0.1 -P 3306 -u letstalk letstalk -puwletstalk
```

#### First time database setup
After starting the server (see below):
```
mysql -h 127.0.0.1 -P 3306 -u letstalk -puwletstalk
msql > CREATE DATABASE letstalk DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;
```

### MAC ONLY

#### Start new docker vm
```
docker-machine start default
```

#### Setup shell environment for docker
Run this in each shell you want to use docker commands from.
```
eval $(docker-machine env)
```

#### Setup port forwarding to your docker container

```
VBoxManage modifyvm "default" --natpf1 "hive,tcp,,8000,,80"
```

Restart machine with updated options.
```
docker-machine restart
```

### LINUX
Things should be good to go

## Installation
See `server/` for server specific installation instruction. See `letstalk/` for client installation instructions.

## Quickstart
Build a docker container and launch the container. Note this will rebuild the server on each file change.
```
./dev.sh
```
NOTE: because of a bug you might have to run `dep ensure` ON YOUR LOCAL MACHINE
since in development mode, the downloaded dependencies will get clobbered.

If you are working with frontend code, you also need to start a separate process
to build the javascript assets so they can be served. You can do this by going to
`server/web` and running:

```
yarn dev
```

## Infrastructure

## Build and startup server on ec2
Run the following command on ec2 server.
```
./prod.sh
```
