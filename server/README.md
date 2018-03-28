# LetsTalk

This is the main server for the UW Let's Talk app.

## Go Development Setup

Get the repository
```
cd $GOPATH/src
git clone git@github.com:andrew749/letstalk.git
```

`$SERVER` will now refer to `$GOPATH/src/letstalk/server`.

Follow the latest instruction for installing go dep at https://golang.github.io/dep/docs/installation.html. This is how we do
dependency management for this project.

### Install/Update Go dependencies
```
go get github.com/codegangsta/gin
dep ensure
```

### Adding new dependencies
```
dep ensure -add github.com/pkg/errors
```

## MySQL

### Development setup

Install mysql, and create database `letstalk`.
```
sudo apt-get install mysql-client mysql-server
mysql -u root -p

CREATE USER letstalk IDENTIFIED BY 'uwletstalk';
CREATE DATABASE letstalk;
GRANT ALL PRIVILEGES ON letstalk . * TO letstalk;
```

## Running the server (development)

```
SECRETS_PATH="secrets.json" ./run_local.sh
```
