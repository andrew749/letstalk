# LetsTalk

This is the main server for the Hive app.

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

