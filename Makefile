APP            := wails/poc
TARGET_RELEASE ?= DEVELOP
COUNTRY ?= ALL

VERSION  ?= $(shell git describe --abbrev=0 --tags)
RELEASE  ?= $(shell git log --pretty=oneline | wc -l | tr -d ' ').g$(shell git rev-parse --short HEAD)
ARCH     := x86_64
#PACKAGE  := $(APP)-$(VERSION)-$(RELEASE).$(ARCH).rpm
#TOP      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(shell dirname $(mkfile_path))
TOP := $(mkfile_dir)
TMP      := $(TOP)/tmp
GOBIN  := $(TOP)/bin
export GOBIN

GOPATH ?= ${HOME}/go
NEWDB ?= "notaname"
PWDD := ${CURDIR}

GOPROXY ?= ""
# Inject the version in a global variable
LDFLAGS  := "-X 'main.version=$(VERSION)-$(RELEASE)'"

GO_COMMONS_BRANCH ?= develop

UNAME_S := $(shell uname -s)
GREP := $(shell which grep)
SUDO := $(shell which sudo)

.PHONY: help build dev\
	goversion

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

build: goversion ##                        Build apps
	$(info) ">>> Install app $(VERSION)-$(RELEASE) ..."
	wails build -tags webkit2_41

dev: goversion ##                        Dev launch
	$(info) ">>> Launch app $(VERSION)-$(RELEASE) ..."
	wails dev -tags webkit2_41

info := @printf "\033[32;01m%s\033[0m\n"
