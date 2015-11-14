build-clean:
	rm -Rf $(PREFIX)/bin/*

extension = $(patsubst windows,.exe,$(filter windows,$(1)))

# Cross builder helper
define gocross
	GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 go build \
		-o $(PREFIX)/bin/docker-machine-driver-sakuracloud-$(1)-$(2)/$(patsubst cmd/%.go,%,$3)$(call extension,$(GOOS)) \
		-a $(VERBOSE_GO) -tags "static_build netgo $(BUILDTAGS)" -installsuffix netgo \
		-ldflags "$(GO_LDFLAGS) -extldflags -static" $(GO_GCFLAGS) $(3);
endef

# XXX building with -a fails in debug (with -N -l) ????

# Independent targets for every bin
$(PREFIX)/bin/%: ./cmd/%.go $(shell find . -type f -name '*.go')
	$(GO) build -o $@$(call extension,$(GOOS)) $(VERBOSE_GO) -tags "$(BUILDTAGS)" -ldflags "$(GO_LDFLAGS)" $(GO_GCFLAGS) $<

# Cross-compilation targets
build-x-%: ./cmd/%.go $(shell find . -type f -name '*.go')
	$(foreach GOARCH,$(TARGET_ARCH),$(foreach GOOS,$(TARGET_OS),$(call gocross,$(GOOS),$(GOARCH),$<)))

# Build just driver
build-driver: $(PREFIX)/bin/docker-machine-driver-sakuracloud

# Overall cross-build
build-x: $(patsubst ./cmd/%.go,build-x-%,$(filter-out %_test.go, $(wildcard ./cmd/*.go)))