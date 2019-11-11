# run from repository root

# Example:
#   make build
#   make build-unix
#   make clean
#   make docker-build

BUILD_DIR ?= ./build
APP_NAME ?= github.com/jaylee630/Hobbit

BUILD_DATE ?= $(shell date +%Y-%m-%d.%H:%M:%S)
GO_VERSION ?= $(shell go version | cut -d ' ' -f 3)
BUILD_VERSION ?= $(shell git rev-parse --short HEAD || echo "GitNotFound")

BUILD_PKG_PATH_PREFIX ?= github.com/jaylee630/Hobbit/ctlmain
LD_FLAGS := "-X $(BUILD_PKG_PATH_PREFIX).buildDate=$(BUILD_DATE) \
	-X $(BUILD_PKG_PATH_PREFIX).buildVersion=$(BUILD_VERSION) \
	-X $(BUILD_PKG_PATH_PREFIX).buildGoVersion=$(GO_VERSION)"

SOURCES := $(shell find . -name '*.go' -not -wholename './vendor/*' -not -wholename '*.pb.go' | xargs echo)


.PHONY: build

all: format test run

build:
	@go build -tags=jsoniter -v -ldflags=$(LD_FLAGS) -o $(BUILD_DIR)/$(APP_NAME)
	@cp -arv ./etc $(BUILD_DIR)/
	$(BUILD_DIR)/$(APP_NAME) -version

test:
	@go test -cover=true  ./...

run: build
	@$(BUILD_DIR)/$(APP_NAME)

format:
	@goimports -w $(SOURCES)
	@echo "Done."

clean:
	rm -rf $(BUILD_DIR)
	rm -f ./github.com/jaylee630/Hobbit
