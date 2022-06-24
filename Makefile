.PHONY: test build

build:
	go build -v main.go

test:
	go test -cover `go list ./... | grep -v 'test'`
