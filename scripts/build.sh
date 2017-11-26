#!/bin/bash

set -e

mkdir -p bin/ 2>/dev/null

for GOOS in $OS; do
    for GOARCH in $ARCH; do
        arch="$GOOS-$GOARCH"
        binary="docker-machine-driver-sakuracloud"
        if [ "$GOOS" = "windows" ]; then
          binary="${binary}.exe"
        fi
        echo "Building $binary $arch"
        GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
            go build \
                -ldflags "$BUILD_LDFLAGS" \
                -o bin/$binary \
                cmd/docker-machine-driver-sakuracloud.go
        if [ -n "$ARCHIVE" ]; then
            (cd bin/; zip -r "docker-machine-driver-sakuracloud_$arch.zip" $binary)
            rm -f bin/$binary
        fi
    done
done
