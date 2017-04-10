FROM golang:1.8
MAINTAINER Kazumichi Yamamoto <yamamoto.febc@gmail.com>

RUN go get  github.com/golang/lint/golint \
            github.com/mattn/goveralls \
            golang.org/x/tools/cover \
            github.com/Azure/go-ansiterm \
            github.com/docker/docker/pkg/system 
            

ENV USER root
WORKDIR /go/src/github.com/yamamoto-febc/docker-machine-sakuracloud

COPY . /go/src/github.com/yamamoto-febc/docker-machine-sakuracloud
