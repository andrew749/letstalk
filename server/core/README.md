# LetsTalk

This is the main service for the UW Let's Talk app.

## Development Setup

```
cd $GOPATH/src
git clone git@github.com:andrew749/letstalk.git
```

`$ROOT` will now refer to `$GOPATH/src/letstalk/server`.

Install Go dependencies:
```
cd $ROOT
go get ./...
```

## MySQL

### Setup

Install mysql, and create database `uwletstalk`.
Linux:
```
sudo apt-get install mysql-client mysql-server
mysql -u root -p
```

```
CREATE USER letstalk IDENTIFIED BY 'uwletstalk';
CREATE DATABASE letstalk;
GRANT ALL PRIVILEGES ON letstalk . * TO letstalk;
```

### Rebuild database

```
cd $ROOT/db
mysql -u letstalk -puwletstalk letstalk < letstalk.sql
```

