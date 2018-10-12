FROM golang:1.9
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN  apt-get update && apt-get -y install bash git make zip && apt-get clean && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
RUN go get -u golang.org/x/lint/golint

ENV USER root
WORKDIR /go/src/github.com/sacloud/docker-machine-sakuracloud
COPY . /go/src/github.com/sacloud/docker-machine-sakuracloud
