GO    := GO15VENDOREXPERIMENT=1 go
pkgs   = $(shell $(GO) list ./... | grep -v /vendor/)
VERSION=$(git describe --tags)                                                   
BUILD=$(date +%FT%T%z)

PREFIX=$(shell pwd)
BIN_DIR=$(shell pwd)


all: format build test

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

test:
	@echo ">> running tests"
	@./runtests.sh

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

build: promu
	@echo ">> building binaries"
	@$(GO) build -ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

tarball: promu
	@echo ">> building release tarball"
	@$(GO) tarball


promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/prometheus/promu


.PHONY: all style format build test vet tarball docker promu
