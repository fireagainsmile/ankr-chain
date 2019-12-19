BUILD_TAGS?='ankrchain'
OUTPUT?=build
NODE_NAME=ankrchain
COMPILER_NAME=contract-compiler

NODE_VERSION=1.0.0
COMPILER_VERSION=1.0.0
CLI_VERSION=1.0.0
LAS_VERSION=1.0.0
NODE_RUNMODE=RunModeUnk
BUILD_FLAGS_NODE = -ldflags "-X github.com/Ankr-network/ankr-chain/version.NodeVersion='git describe --tags $(git rev-list --tags --max-count=1)' -X github.com/Ankr-network/ankr-chain/version.GitCommit=`git rev-parse --short=8 HEAD` X github.com/Ankr-network/ankr-chain/common.RM=${NODE_RUNMODE}"
BUILD_FLAGS_COMPILER = -ldflags "-X github.com/Ankr-network/ankr-chain/version.CompilerVersion='git describe --tags $(git rev-list --tags --max-count=1)' -X github.com/Ankr-network/ankr-chain/version.GitCommit=`git rev-parse --short=8 HEAD`"
BUILD_FLAGS_COMPILER = -ldflags "-X github.com/Ankr-network/ankr-chain/version.CliVersion='git describe --tags $(git rev-list --tags --max-count=1)' -X github.com/Ankr-network/ankr-chain/version.GitCommit=`git rev-parse --short=8 HEAD`"
BUILD_FLAGS_COMPILER = -ldflags "-X github.com/Ankr-network/ankr-chain/version.LasVersion='git describe --tags $(git rev-list --tags --max-count=1)' -X github.com/Ankr-network/ankr-chain/version.GitCommit=`git rev-parse --short=8 HEAD`"

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
    @echo "build ankrchain node image of $(0) for runing mode ${NODE_RUNMODE}"
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(BUILD_FLAGS_NODE) -tags $(BUILD_TAGS) -o $(OUTPUT)/${NODE_NAME}-$(1)-$(2)/$(3) ./main.go
    @echo "build all tools"
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(BUILD_FLAGS_COMPILER) -o ${OUTPUT}/${NODE_NAME}-$(1)-$(2)/$(4) ./tool/compiler/main.go
    CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build  -o ${OUTPUT}/${NODE_NAME}-$(1)-$(2)/$(5) ./tool/cli/main.go
    GGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build  -o ${OUTPUT}/${NODE_NAME}-$(1)-$(2)/$(6) ./service/las/main.go
endef

windows:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,windows,amd64,$(NODE_NAME).exe,$(COMPILER_NAME).exe,$(NODE_NAME)-cli.exe,$(NODE_NAME)-las.exe)

linux:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,linux,amd64,$(NODE_NAME),$(COMPILER_NAME),$(NODE_NAME)-cli,$(NODE_NAME)-las)

darwin:
	@echo "Currency OS:"${PLATFORM}
	$(call build_target,darwin,amd64,$(NODE_NAME),$(COMPILER_NAME),$(NODE_NAME)-cli,$(NODE_NAME)-las)

fmt:
	@go fmt ./...

lint:
	@echo "--> Running linter"
	@golangci-lint run

.PHONY : clean
clean :
	-rm -rf ./build

