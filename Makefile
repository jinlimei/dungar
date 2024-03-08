PROJECT_NAME := "prototype"
PKG := "gitlab.int.magneato.site/dungar/${PROJECT_NAME}"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
GOPATH := $(shell go env GOPATH)
CC := clang-15
CXX := clang-15

export DCWD=$(shell pwd)

.PHONY: all lint dep build clean test coverage coverage_html


all: build

lint: dep ## Lint all the files
	$(GOPATH)/bin/golint -set_exit_status ${PKG_LIST}

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

dep: ## Get the dependencies
	@go mod download
	@go mod verify

build: dep ## Build the binary file
	@go build -o bin/dungar -compiler gc -ldflags="-s -w" cmd/dungar/*.go

cli_test: ## Test on the CLI
	@docker-compose up --force-recreate --remove-orphans --build

clean: ## Remove previous build
	@rm -f bin/dungar

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
