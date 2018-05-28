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

#### How to see mysql tables

Run
```
docker exec -it {DOCKER_CONTAINER_ID} bash
```

Use `docker ps` to list containers

In the container run the following to see mysql database
```
mysql -u letstalk -p
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

## Infrastructure

## Build and startup server on ec2
Run the following command on ec2 server.
```
./prod.sh
```

## Push new image to docker registry (COMING)
Note you need to `docker login` with credentials associated with my docker hub account.
```
docker build . -t andrewcodispoti/hive
docker push andrewcodispoti/hive
```

## Push new files to server
The following script provides an easy way to push files to ec2 in the event we dont want to keep a git repo on the server.
```python infra/push_files.py FILES --destination="~" --username=andrew --private_key=~/.ssh/id_rsa```
