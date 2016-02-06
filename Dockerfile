FROM golang:1.6.0

RUN go get  github.com/golang/lint/golint \
            github.com/mattn/goveralls \
            golang.org/x/tools/cover \
            github.com/tools/godep \
            github.com/aktau/github-release \
            github.com/Azure/go-ansiterm

RUN apt-get update
RUN apt-get install -y zip unzip

ENV USER root
WORKDIR /go/src/github.com/yamamoto-febc/docker-machine-sakuracloud

ADD . /go/src/github.com/yamamoto-febc/docker-machine-sakuracloud
