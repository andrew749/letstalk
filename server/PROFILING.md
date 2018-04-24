# Profiling

## Prerequisites

### INSTALL PPROF
`go get github.com/google/pprof`

### INSTALL GRAPHVIZ
`brew install graphviz`

## How to profile
Run the server with the flag `profiling`

Go to the endpoint `/debug/pprof/profile`. This starts the cpu profile. Now make a bunch of requests to profile.

Eventually the request will finish and now you can visualize this in the browser using `$GOPATH/bin/pprof -http localhost:8080 <PROFILE_FILE>`
