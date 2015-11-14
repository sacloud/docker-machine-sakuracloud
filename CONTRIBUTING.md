# Contributing to docker-machine-sakuracloud

([日本語版](CONTRIBUTING.ja.md))

Want to hack on docker-machine-sakuracloud? Awesome! 
Pull requests are always welcome.

Here are instructions to get you started.


# Building

The requirements to build docker-machine-sakuracloud are:

1. A running instance of Docker or a Golang 1.5 development environment
2. The `bash` shell
3. [Make](https://www.gnu.org/software/make/)

## Build using Docker containers

To build the `docker-machine-sakuracloud` binary using containers, simply run:

    $ export USE_CONTAINER=true
    $ make build

## Local Go 1.5 development environment

Make sure the source code directory is under a correct directory structure to use Go 1.5 vendoring;
example of cloning and preparing the correct environment `GOPATH`:

```
    $ mkdir docker-machine-sakuracloud
    $ cd docker-machine-sakuracloud
    $ export GOPATH="$PWD"
    $ go get github.com/yamamoto-febc/docker-machine-sakuracloud
    $ cd src/github.com/yamamoto-febc/docker-machine-sakuracloud
    $ go get github.com/tools/godep
    $ $GOPATH/bin/godep get
```

At this point, simply run:

    $ make build

## Built binary

After the build is complete a `bin/docker-machine-driver-sakuracloud` binary will be created.

You may call:

    $ make clean

to clean-up build results.

### build targets

Build for all supported oses and architectures (binaries will be in the `bin` project subfolder):

    make cross

Build for a specific list of oses and architectures:

    TARGET_OS=linux TARGET_ARCH="amd64 arm" make cross

You can further control build options through the following environment variables:

    DEBUG=true # enable debug build
    STATIC=true # build static (note: when cross-compiling, the build is always static)
    VERBOSE=true # verbose output
    PREFIX=folder # put binaries in another folder (not the default `./bin`)

Scrub build results:

    make build-clean

### Save and restore dependencies

    make dep-save
    make dep-restore

