GIT_COMMIT = $(shell git rev-parse HEAD)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(BRANCH),HEAD)
	VERSION := $(shell git describe --tags --abbrev=0)
else
	VERSION := $(BRANCH)
endif

BUILD_DATE = $(shell date '+%Y-%m-%d-%H:%M:%S')
HLOADER_VERSION= github.com/cloud-bulldozer/go-commons/version

BIN_DIR = bin
BIN_NAME = hloader
BIN_PATH = $(BIN_DIR)/$(BIN_NAME)
SOURCES = $(shell find . -type f -name "*.go")
CGO = 0

.PHONY: build lint clean

all: lint build container-build

build: $(BIN_PATH)

$(BIN_PATH): $(SOURCES)
	GOARCH=$(shell go env GOARCH) CGO_ENABLED=$(CGO) go build -v -ldflags "-X $(HLOADER_VERSION).GitCommit=$(GIT_COMMIT) -X $(HLOADER_VERSION).Version=$(VERSION) -X $(HLOADER_VERSION).BuildDate=$(BUILD_DATE)" -o $(BIN_PATH) cmd/hloader.go

clean:
	rm -Rf $(BIN_DIR)

lint:
	golangci-lint run
