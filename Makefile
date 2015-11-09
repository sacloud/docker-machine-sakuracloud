default: build

bin/docker-machine-driver-sakuracloud:
	go build -i -o ./bin/docker-machine-driver-sakuracloud ./bin

build: clean bin/docker-machine-driver-sakuracloud

clean:
	$(RM) bin/docker-machine-driver-sakuracloud

install: bin/docker-machine-driver-sakuracloud
	cp -f ./bin/docker-machine-driver-sakuracloud $(GOPATH)/bin/
