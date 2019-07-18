BUILD_TAGS?='ankrchain'
OUTPUT?=build/ankrchain
BUILD_FLAGS = -ldflags "-X github.com/Ankr-network/ankr-chain/version.GitCommit=`git rev-parse --short=8 HEAD`"

export GO111MODULE=on

all: build install

build:
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -tags $(BUILD_TAGS) -o $(OUTPUT) ./main.go

install:
	CGO_ENABLED=0 go install  $(BUILD_FLAGS) -tags $(BUILD_TAGS) ./main.go

fmt:
	@go fmt ./...

lint:
	@echo "--> Running linter"
	@golangci-lint run

.PHONY: check build install fmt lint

