BINARY_NAME=acli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS  = -s -w \
	-X github.com/chinmaymk/acli/cmd/acli.version=$(VERSION) \
	-X github.com/chinmaymk/acli/cmd/acli.commit=$(COMMIT) \
	-X github.com/chinmaymk/acli/cmd/acli.date=$(DATE)

.PHONY: build clean test lint install all

build:
	go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME) .

install:
	go install -ldflags '$(LDFLAGS)' .

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/ dist/

# Cross-compile for common platforms
all: clean
	GOOS=darwin  GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o bin/$(BINARY_NAME)-windows-amd64.exe .
