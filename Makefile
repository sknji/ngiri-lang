VERSION=0.9

build:
	@go build -o ngiri cli/cli.go

fmt:
	@go fmt ./...

tests:
	@go test ./...

.PHONY: build clean

clean:
	@rm -r ngiri

