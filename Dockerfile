FROM ubuntu:16.04
MAINTAINER Andrew Codispoti

RUN apt-get update -y
RUN apt-get install -y postgresql redis-server scala

RUN apt-get install -y curl

ENV SCALA_VERSION 2.12.4
ENV SBT_VERSION 1.0.2

# Install Scala
## Piping curl directly in tar
RUN \
  curl -fsL https://downloads.typesafe.com/scala/$SCALA_VERSION/scala-$SCALA_VERSION.tgz | tar xfz - -C /root/ && \
  echo >> /root/.bashrc && \
  echo "export PATH=~/scala-$SCALA_VERSION/bin:$PATH" >> /root/.bashrc

# Install sbt
RUN \
  curl -L -o sbt-$SBT_VERSION.deb https://dl.bintray.com/sbt/debian/sbt-$SBT_VERSION.deb && \
  dpkg -i sbt-$SBT_VERSION.deb && \
  rm sbt-$SBT_VERSION.deb && \
  apt-get update && \
  apt-get install sbt && \
  sbt sbtVersion

# Add server/app using docker run so we have a mounted filesystem

# set the working directory
WORKDIR /home/app
CMD sbt run messengerservice/
EXPOSE 8080
