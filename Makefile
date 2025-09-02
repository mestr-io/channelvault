APP            := wails/poc
TARGET_RELEASE ?= DEVELOP
COUNTRY ?= ALL

#VERSION  ?= $(shell git describe --abbrev=0 --tags) #when we start tagging versions this may be useful
VERSION  ?= "0.1"
RELEASE  ?= $(shell git log --pretty=oneline | wc -l | tr -d ' ').g$(shell git rev-parse --short HEAD)
ARCH     := x86_64
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell dirname $(mkfile_path))
TOP := $(mkfile_dir)
TMP := $(TOP)/tmp
GOBIN  := $(TOP)/bin
export GOBIN

GOPATH ?= ${HOME}/go
GOPROXY ?= ""
# Inject the version in a global variable
LDFLAGS  := "-X 'main.version=$(VERSION)-$(RELEASE)'"

.PHONY: help build dev compile goversion fmt check-fmt test cover clean testh coverh

all: help

help:##                         Show this help.
	@echo
	@echo "usage: make <command>"
	@echo
	@echo "commands:"
	@echo
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -nr 's/(.*):.*##(.*)/\1 \2/p'

goversion:
	$(info) ">>> Which go version ..."
	@go version

compile: goversion ##                      Just compile
	$(info) ">>> Compile app $(VERSION)-$(RELEASE) ..."
	@go install -v

build: goversion ##                        Build apps
	$(info) ">>> Install app $(VERSION)-$(RELEASE) ..."
	@wails build -tags webkit2_41

dev: goversion ##                          Dev launch
	$(info) ">>> Launch app $(VERSION)-$(RELEASE) ..."
	wails dev -tags webkit2_41

test: goversion ##                         Launch app tests
	$(info) ">>> Launch tests for app $(VERSION)-$(RELEASE) ..."
	@go test -cover -count=1 -v -tags=dev,webkit2_41 ./...

cover: goversion ##                        Launch app tests just cover report
	$(info) ">>> Launch tests for app $(VERSION)-$(RELEASE) ..."
	@go test -cover -count=1 -tags=dev,webkit2_41 ./...

testh: goversion ##                        Launch headless app tests
	$(info) ">>> Launch headless tests for app $(VERSION)-$(RELEASE) ..."
	DISPLAY=:99 xvfb-run -a go test -cover -count=1 -v -tags=dev,headless,webkit2_41 ./...

coverh: goversion ##                       Launch headless app tests just cover report
	$(info) ">>> Launch headless tests for app $(VERSION)-$(RELEASE) ..."
	DISPLAY=:99 xvfb-run -a go test -cover -count=1 -tags=dev,headless,webkit2_41 ./...

jsdeps: ##                       Refresh nodejs deps
	cd $(TOP)/frontend && npm update && npm install

lint: ##                         launch go lint
	$(info) ">>> Launch headless tests for app $(VERSION)-$(RELEASE) ..."
	@golangci-lint run

clean: ##                        Clean build
	rm -Rf $(TOP)/build/bin
	rm -Rf $(TOP)/frontend/dist
	rm -Rf $(TOP)/frontend/node_modules
	mkdir -p $(TOP)/frontend/dist
	touch $(TOP)/frontend/dist/index.html

fmt:
	$(info) ">>> Format JS/TS code $(VERSION)-$(RELEASE) ..."
	cd $(TOP)/frontend && \
	$(TOP)/frontend/node_modules/@biomejs/biome/bin/biome format --write

fmt-check:	
	$(info) ">>> Check format on JS/TS code $(VERSION)-$(RELEASE) ..."
	cd $(TOP)/frontend && \
	$(TOP)/frontend/node_modules/@biomejs/biome/bin/biome format

info := @printf "\033[32;01m%s\033[0m\n"
