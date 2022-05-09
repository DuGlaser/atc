.PHONY: test build

test:
	go test -cover `go list ./... | grep -v 'test'`

build:
	go build -v main.go
