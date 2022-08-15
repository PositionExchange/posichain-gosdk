SHELL := /bin/bash
version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}@posichain.org

flags := -gcflags="all=-N -l -c 2"
ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}
cli := ./dist/psc
uname := $(shell uname)

env := GO111MODULE=on

DIR := ${CURDIR}
export CGO_LDFLAGS=-L$(DIR)/dist/lib -Wl,-rpath -Wl,\$ORIGIN/lib

all:
	source $(shell go env GOPATH)/src/github.com/PositionExchange/posichain-gosdk/scripts/setup_bls_build_flags.sh && $(env) go build -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) psc

static:
	make -C $(shell go env GOPATH)/src/github.com/PositionExchange/mcl
	make -C $(shell go env GOPATH)/src/github.com/PositionExchange/bls minimised_static BLS_SWAP_G=1
	source $(shell go env GOPATH)/src/github.com/PositionExchange/posichain-gosdk/scripts/setup_bls_build_flags.sh && $(env) go build -o $(cli) -ldflags="$(ldflags) -w -extldflags \"-static\"" cmd/main.go
	cp $(cli) psc

debug:
	source $(shell go env GOPATH)/src/github.com/PositionExchange/posichain-gosdk/scripts/setup_bls_build_flags.sh && $(env) go build $(flags) -o $(cli) -ldflags="$(ldflags)" cmd/main.go
	cp $(cli) psc

install:all
	cp $(cli) ~/.local/bin

run-tests: test-rpc test-key;

test-key:
	go test ./pkg/keys -cover -v

test-rpc:
	go test ./pkg/rpc -cover -v

.PHONY:clean run-tests

clean:
	@rm -f $(cli)
	@rm -rf ./dist