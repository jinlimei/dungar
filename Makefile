PKG := "gitlab.int.magneato.site/dungar/dungar"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOPATH := $(shell go env GOPATH)
CC := clang-22
CXX := clang-22

export DCWD=$(shell pwd)
export GOAMD64=v3

.PHONY: all dep build fix vet ci_lint lint slim_lint test race msan coverage coverage_html cli_test clean help

all: build

dep: ## Get the dependencies
	@go mod download
	@go mod verify

build: dep ## Build the binary file
	@go generate ./...
	@go build -o bin/dungar -compiler gc -ldflags="-s -w" cmd/dungar/*.go

fix: dep
	@go generate ./...
	@go fix ${PKG_LIST}
	@go vet -fix ${PKG_LIST}

vet: dep ## go vet
	@go vet ${PKG_LIST}

ci_lint: dep ## golangci-lint
	@golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --output.code-climate.path=gl-code-quality-report.json

lint: dep ## golangci-lint
	@golangci-lint run --max-issues-per-linter 0 --max-same-issues 0

slim_lint: dep ## golangci-lint
	@golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --new

test: ## Run tests
	@go test -parallel 8 ${PKG_LIST}

race: ## Run race detection
	@go test -race -parallel 8 ${PKG_LIST}

msan: ## Run data race detector
	@go test -msan -parallel 8 ${PKG_LIST}

coverage: ## Generate global code coverage report
	@bash ./tools/coverage.sh;

coverage_html:
	@bash ./tools/coverage.sh html;

cli_test: ## Test on the CLI
	@docker compose up --force-recreate --remove-orphans --build

clean: ## Remove previous build
	@rm -f bin/dungar

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
