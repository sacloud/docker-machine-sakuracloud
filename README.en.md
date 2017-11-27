# Docker Machine SAKURA CLOUD Driver

This is a plugin for [Docker Machine](https://docs.docker.com/machine/) allowing
to create Docker hosts on [SAKURA CLOUD](http://cloud.sakura.ad.jp)

([日本語版](README.md))

[![Build Status](https://travis-ci.org/sacloud/docker-machine-sakuracloud.svg?branch=master)](https://travis-ci.org/sacloud/docker-machine-sakuracloud)

## Quick Start: Running with Docker

We released docker image that bundle docker-machine and docker-machine-driver-sakuracloud.

[sacloud/docker-machine](https://hub.docker.com/r/sacloud/docker-machine/)

If you don't have docker, try [Install Driver](#install-driver).

```bash
docker run [docker run options] sacloud/docker-machine [docker-machine options] <machine-name>
```

Example: 

```bash
docker run -it --rm -e MACHINE_STORAGE_PATH=$HOME/.docker/machine \
                    -e SAKURACLOUD_ACCESS_TOKEN=<Your Access Token> \
                    -e SAKURACLOUD_ACCESS_TOKEN_SECRET=<Your Access Secret> \
                    -v $HOME/.docker:$HOME/.docker \
                    sacloud/docker-machine create -d sakuracloud sakura-dev
```

## Install Driver

### macOS(HomeBrew)

    $ brew tap sacloud/docker-machine-sakuracloud
    $ brew install docker-machine-sakuracloud


### Manual Install

download the binary `docker-machine-driver-sakuracloud`
and  make it available by `$PATH`, for example by putting it to `/usr/local/bin/`

The latest version of `docker-machine-driver-sakuracloud` binary is available on
the ["Releases"](https://github.com/sacloud/docker-machine-sakuracloud/releases/latest) page.

## Usage
Official documentation for Docker Machine [is available here](https://docs.docker.com/machine/).

To create a virtual machine on `SAKURA CLOUD` for Docker purposes just run this command:

```
$ docker-machine create --driver=sakuracloud \
    --sakuracloud-access-token=<YOUR TOKEN> \
    --sakuracloud-access-token-secret=<YOUR TOKEN SECRET> \
    sakura-dev
```

Options:

 - `--sakuracloud-access-token`: **required** Your personal access token for the SAKURA CLOUD API.
 - `--sakuracloud-access-token-secret`: **required** Your personal access token secret for the SAKURA CLOUD API.
 - `--sakuracloud-zone`: Zone [`is1a` / `is1b` / `tk1a`]
 - `--sakuracloud-os-type`: OS type [`rancheros` / `centos` / `ubuntu`]
 - `--sakuracloud-core`: Number of CPU-core
 - `--sakuracloud-memory`: Size of memory (In GB)
 - `--sakuracloud-disk-connection`: Disk connection type(`virtio` or `ide`)
 - `--sakuracloud-disk-plan`: Disk plan(`ssd` / `hdd`)
 - `--sakuracloud-disk-size`: Size of disk(In GB)
 - `--sakuracloud-interface-driver`: Interface driver(`virtio` or `e1000`)
 - `--sakuracloud-password`: Password for Admin user(if empty, use random strings)
 - `--sakuracloud-enable-password-auth` : Enable password auth when connect by SSH
 - `--sakuracloud-packet-filter`: ID of packet filter
 - `--sakuracloud-engine-port` : The number of DockerEngine port.
 - `--sakuracloud-ssh-key` : The path of ssh private key.

Environment variables and default values:

| CLI option                           | Environment variable              | Default                  |
|--------------------------------------|-----------------------------------|--------------------------|
| `--sakuracloud-access-token`         | `SAKURACLOUD_ACCESS_TOKEN`        | -                        |
| `--sakuracloud-access-token-secret`  | `SAKURACLOUD_ACCESS_TOKEN_SECRET` | -                        |
| `--sakuracloud-zone`                 | `SAKURACLOUD_ZONE`                | `is1b`                   |
| `--sakuracloud-os-type`              | `SAKURACLOUD_OS_TYPE`             | `rancheros`              |
| `--sakuracloud-core`                 | `SAKURACLOUD_CORE`                | `1`                      |
| `--sakuracloud-memory`               | `SAKURACLOUD_MEMORY`              | `1`                      |
| `--sakuracloud-disk-connection`      | `SAKURACLOUD_DISK_CONNECTION`     | `virtio`                 |
| `--sakuracloud-disk-plan`            | `SAKURACLOUD_DISK_PLAN`           | `ssd`                    |
| `--sakuracloud-disk-size`            | `SAKURACLOUD_DISK_SIZE`           | `20`                     |
| `--sakuracloud-interface-driver`     | `SAKURACLOUD_INTERFACE_DRIVER`    | `virtio`                 |
| `--sakuracloud-password`             | `SAKURACLOUD_PASSWORD`            | -                        |
| `--sakuracloud-enable-password-auth` | `SAKURACLOUD_ENABLE_PASSWORD_AUTH`| false                    |
| `--sakuracloud-packet-filter`        | `SAKURACLOUD_PACKET_FILTER`       | -                        |
| `--sakuracloud-engine-port`          | `SAKURACLOUD_ENGINE_PORT`         | `2376`                   |
| `--sakuracloud-ssh-key`              | `SAKURACLOUD_SSH_KEY`             | -                        |


## Author

* Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
