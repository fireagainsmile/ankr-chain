BUILD_TAGS?='ankrchain'
OUTPUT?=build
BUILD_FLAGS = -ldflags "-X github.com/tendermint/tendermint/version.GitCommit=`git rev-parse --short=8 HEAD`"
NODE_NAME=ankrchain
COMPILER_NAME=contract-compiler

OUTPUTTOOLDIR=build/tool

export GO111MODULE=on

ifeq ($(OS),Windows_NT)
  PLATFORM="Windows"
else
  ifeq ($(shell uname),Darwin)
    PLATFORM="MacOS"
  else
    PLATFORM="Linux"
  endif
endif

all: windows linux darwin

define build_target
    @echo "build ankrchain node image of $(0)"
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(BUILD_FLAGS) -tags $(BUILD_TAGS) -o $(OUTPUT)/${NODE_NAME}-$(1)-$(2)/$(3) ./main.go
    @echo "build all tools"
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build  -o ${OUTPUT}/${NODE_NAME}-$(1)-$(2)/$(4) ./tool/compiler/main.go
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build  -o ${OUTPUT}/${NODE_NAME}-$(1)-$(2)/$(5) ./tool/cli/main.go
endef

windows:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,windows,amd64,$(NODE_NAME).exe,$(COMPILER_NAME).exe,$(NODE_NAME)-cli.exe)

linux:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,linux,amd64,$(NODE_NAME),$(COMPILER_NAME),$(NODE_NAME)-cli)

darwin:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,darwin,amd64,$(NODE_NAME),$(COMPILER_NAME),$(NODE_NAME)-cli)

fmt:
	@go fmt ./...

lint:
	@echo "--> Running linter"
	@golangci-lint run

.PHONY : clean
clean :
	-rm -rf ./build

