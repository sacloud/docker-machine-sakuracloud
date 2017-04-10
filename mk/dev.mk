test:
	go test $(shell go list ./... | grep -v vendor/)