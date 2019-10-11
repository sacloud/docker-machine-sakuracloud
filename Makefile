TEST?=$$(go list ./...)
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOLINT_TARGETS?=$$(golint github.com/sacloud/docker-machine-sakuracloud | tee /dev/stderr)
CURRENT_VERSION := $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')

BUILD_LDFLAGS = "-s -w \
	  -X `go list ./version`.GitCommit=`git rev-parse --short HEAD 2>/dev/null` \
	  -X `go list ./version`.Version=$(CURRENT_VERSION) "

export GO111MODULE=on

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux shasum

build-darwin: bin/docker-machine-driver-sakuracloud_darwin-386.zip bin/docker-machine-driver-sakuracloud_darwin-amd64.zip

build-windows: bin/docker-machine-driver-sakuracloud_windows-386.zip bin/docker-machine-driver-sakuracloud_windows-amd64.zip

build-linux: bin/docker-machine-driver-sakuracloud_linux-386.zip bin/docker-machine-driver-sakuracloud_linux-amd64.zip

bin/docker-machine-driver-sakuracloud_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/docker-machine-driver-sakuracloud_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/docker-machine-driver-sakuracloud_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/docker-machine-driver-sakuracloud_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/docker-machine-driver-sakuracloud_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/docker-machine-driver-sakuracloud_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

shasum:
	(cd bin/; shasum -a 256 * > docker-machine-driver-sakuracloud_SHA256SUMS)

test: vet
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4 ; \

vet: golint
	go vet ./...

golint:
	test -z "$$(golint ./... | grep -v 'tools/' | grep -v 'vendor/' | grep -v '_string.go' | tee /dev/stderr )"

goimports: fmt
	goimports -l -w $(GOFMT_FILES)

fmt:
	gofmt -s -l -w $(GOFMT_FILES)

docker-test:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"

.PHONY: default test vet testacc fmt fmtcheck
