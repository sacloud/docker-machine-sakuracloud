FROM golang:1.16 AS builder
MAINTAINER Sacloud Authors <sacloud.users@gmail.com>

RUN  apt-get update && apt-get -y install \
        bash \
        git  \
        make \
        zip  \
        bzr  \
      && apt-get clean \
      && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

ADD . /go/src/github.com/sacloud/docker-machine-sakuracloud
WORKDIR /go/src/github.com/sacloud/docker-machine-sakuracloud
ENV CGO_ENABLED 0
RUN make tools build
# ======

FROM alpine:3.13
MAINTAINER Sacloud Authors <sacloud.users@gmail.com>

ADD https://github.com/docker/machine/releases/download/v0.16.2/docker-machine-Linux-x86_64 /usr/local/bin/docker-machine
RUN chmod +x /usr/local/bin/docker-machine
RUN apk add --no-cache --update ca-certificates
COPY --from=builder /go/src/github.com/sacloud/docker-machine-sakuracloud/bin/docker-machine-driver-sakuracloud /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/docker-machine"]
