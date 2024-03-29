.PHONY: build test clean all docs
.DEFAULT_GOAL := all

SHELL := /bin/bash

git_branch := $(shell git rev-parse --abbrev-ref HEAD)

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

update-pace-protos-to-latest-tag:
	@ latest_tag=$$(BUF_BETA_SUPPRESS_WARNINGS=1 buf beta registry tag list buf.build/getstrm/pace --reverse --page-size 1 --format json | jq -r '.results[0].name') && \
	echo "The latest tag for the PACE protos is: $$latest_tag" && \
	BUF_ALPHA_SUPPRESS_WARNINGS=1 buf alpha sdk go-version --module=buf.build/getstrm/pace:$$latest_tag --plugin=buf.build/grpc/go:v1.3.0 \
	| xargs -I% go get buf.build/gen/go/getstrm/pace/grpc/go@% && go mod tidy

update-pace-protos-to-current-git-branch:
	BUF_ALPHA_SUPPRESS_WARNINGS=1 buf alpha sdk go-version --module=buf.build/getstrm/pace:${git_branch} --plugin=buf.build/grpc/go:v1.3.0 \
	| xargs -I% go get buf.build/gen/go/getstrm/pace/grpc/go@% && go mod tidy
