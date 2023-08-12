.PHONY: test build generate/lang_list

build:
	go build -v main.go

test:
	go test -cover `go list ./... | grep -v 'test'`

generate/lang_list:
	@go run ./scripts/generate_lang_json.go $(TARGET) | jq .


