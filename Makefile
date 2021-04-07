GO_FILES  ?=$(shell find . -name '*.go')
TEST      ?=$$(go list ./...)

BUILD_LDFLAGS = "-s -w -X `go list ./version`.GitCommit=`git rev-parse --short HEAD 2>/dev/null`"

export GO111MODULE=on

.PHONY: default
default: test

.PHONY: tools
tools:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get golang.org/x/tools/cmd/stringer
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.37.0/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.37.0

.PHONY: clean
clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: build
build: bin/docker-machine-driver-sakuracloud

bin/docker-machine-driver-sakuracloud: $(GO_FILES)
	go build -ldflags $(BUILD_LDFLAGS) -o bin/docker-machine-driver-sakuracloud cmd/docker-machine-driver-sakuracloud.go

.PHONY: test
test:
	TF_ACC= go test $(TEST) $(TESTARGS) -timeout=30s -parallel=4

.PHONY: lint
lint:
	golangci-lint run --modules-download-mode=readonly ./...

.PHONY: goimports
goimports: fmt
	goimports -l -w .

.PHONY: fmt
fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

