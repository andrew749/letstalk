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

## Quickstart
Build a docker container and launch the container. Note this will rebuild the server on each file change.
```
./dev.sh
```

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
