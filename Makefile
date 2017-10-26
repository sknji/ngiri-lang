VERSION=0.9

build:
	@go build -o monkey cli/cli.go

fmt:
	@go fmt ./...

.PHONY: build clean

clean:
	@rm -r monkey
