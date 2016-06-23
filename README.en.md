# Docker Machine SAKURA CLOUD Driver

This is a plugin for [Docker Machine](https://docs.docker.com/machine/) allowing
to create Docker hosts on [SAKURA CLOUD](http://cloud.sakura.ad.jp)

([日本語版](README.md))

[![Build Status](https://travis-ci.org/yamamoto-febc/docker-machine-sakuracloud.svg?branch=master)](https://travis-ci.org/yamamoto-febc/docker-machine-sakuracloud)

## Requirements
* [Docker Machine](https://docs.docker.com/machine/) 0.5.1+ (is bundled to
  [Docker Toolbox](https://www.docker.com/docker-toolbox) 1.9.1+)

## Tested Operationg System
* OSX 10.9+
* Windows 10

## Installation

#### Install for Windows

Download the installer from [here](https://github.com/yamamoto-febc/docker-machine-sakuracloud/releases/download/v0.0.12/DockerMachineSakuracloudSetup.exe)
 and run it.

#### Install for OSX:

download the binary `docker-machine-driver-sakuracloud`
and  make it available by `$PATH`, for example by putting it to `/usr/local/bin/`

```console
curl -L https://github.com/yamamoto-febc/docker-machine-sakuracloud/releases/download/v0.0.12/docker-machine-driver-sakuracloud-`uname -s`-`uname -m` >/usr/local/bin/docker-machine-driver-sakuracloud && \
  chmod +x /usr/local/bin/docker-machine/docker-machine-driver-sakuracloud
```

The latest version of `docker-machine-driver-sakuracloud` binary is available on
the ["Releases"](https://github.com/yamamoto-febc/docker-machine-sakuracloud/releases/latest) page.

## Usage
Official documentation for Docker Machine [is available here](https://docs.docker.com/machine/).

To create a virtual machine on `SAKURA CLOUD` for Docker purposes just run this command:

```
$ docker-machine create --driver=sakuracloud \
    --sakuracloud-access-token=[YOUR TOKEN] \
    --sakuracloud-access-token-secret=[YOUR TOKEN SECRET] \
    sakura-dev
```

Options:

 - `--sakuracloud-access-token`: **required** Your personal access token for the SAKURA CLOUD API.
 - `--sakuracloud-access-token-secret`: **required** Your personal access token secret for the SAKURA CLOUD API.
 - `--sakuracloud-core`: The number of CPU cores.
 - `--sakuracloud-connected-switch`: The ID of SAKURA CLOUD switch or router.
 - `--sakuracloud-disk-connection`: The type of disk connection (`virtio` or `ide`).
 - `--sakuracloud-disk-name`: The name of SAKURA CLOUD disk.
 - `--sakuracloud-disk-plan`: The plan of SAKURA CLOUD disk plan (HDD:`2` or SSD:`4`).
 - `--sakuracloud-disk-size`: The size of disk for the SAKURA CLOUD server(in MB).
 - `--sakuracloud-dns-zone` : The domain name for SAKURACLOUD DNS.
 - `--sakuracloud-memory-size`: The size of Memory(in GB).
 - `--sakuracloud-private-ip-only`: The flag of to use private IP only(use eth1 only).
 - `--sakuracloud-private-ip`: The IP address for eth1.
 - `--sakuracloud-private-ip-subnet-mask`: The subnet mask for eth1.
 - `--sakuracloud-gateway`: The default gateway ip address.
 - `--sakuracloud-region`: The resion to create the server in.
 - `--sakuracloud-group`: The @group tag.
 - `--sakuracloud-gslb`: The Name of GSLB.
 - `--sakuracloud-auto-reboot`: The @auto-reboot tag.
 - `--sakuracloud-ignore-virtio-net`: The flag of to not set @virtio-net-pci tag.
 - `--sakuracloud-packet-filter`: The Packet Filter ID or Name for eth0(shared).
 - `--sakuracloud-private-packet-filter`: The Packet Filter ID or Name for eth1(private).
 - `--sakuracloud-enable-password-auth` : The flag of enable password auth.
 - `--sakuracloud-engine-port` : The number of DockerEngine port.
 - `--sakuracloud-ssh-key` : The path of ssh private key.

Environment variables and default values:

| CLI option                          | Environment variable              | Default                  |
|-------------------------------------|-----------------------------------|--------------------------|
| `--sakuracloud-access-token`        | `SAKURACLOUD_ACCESS_TOKEN`        | -                        |
| `--sakuracloud-access-token-secret` | `SAKURACLOUD_ACCESS_TOKEN_SECRET` | -                        |
| `--sakuracloud-auto-reboot`         | `SAKURACLOUD_AUTO_REBOOT`        | -                   |
| `--sakuracloud-core`                | `SAKURACLOUD_CORE`                | `1`                   |
| `--sakuracloud-connected-switch`    | `SAKURACLOUD_CONNECTED_SWITCH`     | -                 |
| `--sakuracloud-disk-connection`     | `SAKURACLOUD_DISK_CONNECTION`     | `virtio`                 |
| `--sakuracloud-disk-name`           | `SAKURACLOUD_DISK_NAME`           | `disk001`                |
| `--sakuracloud-disk-plan`           | `SAKURACLOUD_DISK_PLAN`           | `4`                      |
| `--sakuracloud-disk-size`           | `SAKURACLOUD_DISK_SIZE`           | `20480`                  |
| `--sakuracloud-dns-zone`   | `SAKURACLOUD_DNS_ZONE`  | -                 |
| `--sakuracloud-enable-password-auth`   | `SAKURACLOUD_ENABLE_PASSWORD_AUTH`  | false                 |
| `--sakuracloud-engine-port`   | `SAKURACLOUD_ENGINE_PORT`  | `2376`                 |
| `--sakuracloud-gateway`     | `SAKURACLOUD_GATEWAY`     | -                 |
| `--sakuracloud-group`               | `SAKURACLOUD_GROUP`              | -                   |
| `--sakuracloud-gslb`               | `SAKURACLOUD_GSLB`              | -                   |
| `--sakuracloud-ignore-virtio-net`   | `SAKURACLOUD_IGNORE_VIRTIO_NET`  | -                   |
| `--sakuracloud-memory-size`         | `SAKURACLOUD_MEMORY_SIZE`         | `1`                   |
| `--sakuracloud-packet-filter`   | `SAKURACLOUD_PACKET_FILTER`  | -                   |
| `--sakuracloud-private-ip-only`       | `SAKURACLOUD_PRIVATE_IP_ONLY`     | -                 |
| `--sakuracloud-private-ip`       | `SAKURACLOUD_PRIVATE_IP`     | -                 |
| `--sakuracloud-private-ip-subnet-mask`     | `SAKURACLOUD_PRIVATE_IP_SUBNET_MASK`     | `255.255.255.0`          |
| `--sakuracloud-private-packet-filter`   | `SAKURACLOUD_PRIVATE_PACKET_FILTER`  | -                   |
| `--sakuracloud-region`              | `SAKURACLOUD_REGION`              | `is1b`                   |
| `--sakuracloud-ssh-key`   | `SAKURACLOUD_SSH_KEY`  | -                 |


## Author

* Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
