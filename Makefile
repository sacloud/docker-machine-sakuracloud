default: build

bin/docker-machine-driver-sakuracloud:
	go build -i -o ./bin/docker-machine-driver-sakuracloud ./bin

build: clean bin/docker-machine-driver-sakuracloud

clean:
	$(RM) bin/docker-machine-driver-sakuracloud \
	rm -Rf bin/windows_*;	rm -Rf bin/linux_* ; 	rm -Rf bin/darwin_*


install: bin/docker-machine-driver-sakuracloud
	cp -f ./bin/docker-machine-driver-sakuracloud $(GOPATH)/bin/

build-windows-386:
	GOOS=windows GOARCH=386 go build -i -o ./bin/windows_386/docker-machine-driver-sakuracloud.exe ./bin
build-windows-64:
	GOOS=windows GOARCH=amd64 go build -i -o ./bin/windows_amd64/docker-machine-driver-sakuracloud.exe ./bin
build-linux-386:
	GOOS=linux GOARCH=386 go build -i -o ./bin/linux_386/docker-machine-driver-sakuracloud ./bin
build-linux-64:
	GOOS=linux GOARCH=amd64 go build -i -o ./bin/linux_amd64/docker-machine-driver-sakuracloud ./bin
build-darwin-386:
	GOOS=darwin GOARCH=386 go build -i -o ./bin/darwin_386/docker-machine-driver-sakuracloud ./bin
build-darwin-64:
	GOOS=darwin GOARCH=amd64 go build -i -o ./bin/darwin_amd64/docker-machine-driver-sakuracloud ./bin
