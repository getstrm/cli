.PHONY: build test clean all docs
.DEFAULT_GOAL := all

SHELL := /bin/bash

build:
	goreleaser --snapshot --skip-publish --clean

zsh-completion:
	pace completion zsh > "$${fpath[1]}/_pace"

# for a speedier build than with goreleaser
source_files := $(shell find . -name "*.go")

targetVar := pace/pace/pkg/util.RootCommandName

target := dpace

ldflags := -X '${targetVar}=${target}' -X pace/pace/pkg/cmd.Version=local -X pace/pace/pkg/common.GitSha=local -X pace/pace/pkg/common.BuiltOn=local

dist/${target}: ${source_files} Makefile
	go build -ldflags="${ldflags}" -o $@ ./cmd/pace

clean:
	rm -f dist/${target} dist/pace

# Make sure the .env containing all `STRM_TEST_*` variables is present in the ./test directory
# godotenv loads the .env file from that directory when running the tests
test: dist/${target}
	go clean -testcache
	go test ./test -v

all: dist/${target}

dist/pace: ${source_files} Makefile
	go build -o $@ ./cmd/pace

docs: dist/pace
	dist/pace generate-docs
