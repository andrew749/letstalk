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

### Easy database/modelq setup

run `$SERVER/db/update_db.sh`.

### Rebuilding the database

An up-to-date schema of the entire database should be kept under
`$SERVER/db`.

```
cd $SERVER/db
mysql -u letstalk -puwletstalk letstalk < letstalk.sql
```

### Running modelq

Install modelq:
```
go get github.com/mijia/modelq
```

Create a `modelq` user with minimal privileges:
```
mysql -u root -p

CREATE USER modelq;
GRANT REFERENCES ON letstalk . * TO modelq;
```

Modelq automatically generates data structures and ORM functions from
the current local database schema. Generated files are placed
in `$SERVER/data`.

Generate code from schema:
```
cd $SERVER
modelq -db="modelq@/letstalk" -pkg=data -driver=mysql -schema=letstalk -p=8
```
