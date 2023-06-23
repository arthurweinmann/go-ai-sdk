BUILDDIR                 := $(CURDIR)/build


# It's necessary to set this because some environments don't link sh -> bash.
SHELL                             = /usr/bin/env bash

GOPATH                            = $(shell go env GOPATH)
GOBIN                             = $(shell which go)
ARCH                              = $(shell uname -p)

GIT_COMMIT                        = $(shell git rev-parse HEAD)
GIT_SHA                           = $(shell git rev-parse --short HEAD)
GIT_TAG                           = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY                         = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# +++ Check bash installed

DEPTEST=$(shell command -v $(SHELL) 2> /dev/null)
ifeq ($(DEPTEST),)
$(error "We could not find bash")
endif


# +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


# --------------------------------------------------------------------------------
# --------------------------------------------------------------------------------
# --------------------------------------------------------------------------------

.PHONY: test
test: check-generic-dep
	@cd $(CURDIR)/test && go test -v

PHONY:check-generic-dep
check-generic-dep:
	@command -v git >/dev/null 2>&1 || { echo >&2 "git is not installed or not in path"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo >&2 "go is not installed or not in path"; exit 1; }
	@command -v python3 >/dev/null 2>&1 || { echo >&2 "python3 is not installed or not in path"; exit 1; }