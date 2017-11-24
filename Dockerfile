
FROM ubuntu:16.04
MAINTAINER Andrew Codispoti

RUN apt-get update -y
RUN apt-get install -y postgresql redis-server scala

RUN apt-get install -y curl

RUN \
  curl -L -o sbt-$SBT_VERSION.deb http://dl.bintray.com/sbt/debian/sbt-$SBT_VERSION.deb && \
  dpkg -i sbt-$SBT_VERSION.deb && \
  rm sbt-$SBT_VERSION.deb && \
  apt-get update && \
  apt-get install sbt && \
  sbt sbtVersion

ADD server/ /home/app/

WORKDIR /home/app
CMD sbt run messengerservice/
EXPOSE 8080
