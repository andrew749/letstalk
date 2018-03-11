FROM ubuntu:16.04
MAINTAINER Andrew Codispoti

RUN apt-get update -y
RUN apt-get install -y golang git
RUN apt-get install -y curl

# gopath in root
ENV GOPATH /go
ENV PATH="${PATH}:/go/bin"

# gin to run reloading server
RUN go get github.com/codegangsta/gin

# add the source code
ADD ./server /go/src/letstalk/server

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# set the working directory
WORKDIR go/src/letstalk/server
RUN dep ensure

# fetch dependencies
ENV SECRETS_PATH secrets.json
CMD ./run_local.sh
EXPOSE 3000
