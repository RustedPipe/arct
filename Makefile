.phony: build test snapshot

ROOT                := $(PWD)
GO_HTML_COV         := ./coverage.html
GO_TEST_OUTFILE     := ./c.out
GO_DOCKER_IMAGE     := golang:1.16
GO_DOCKER_CONTAINER := arct-container

SHELL     := /bin/bash
PROJECT   := arct
VCS_REF   := `git rev-parse HEAD`
ITERATION := $(shell date -u +%Y-%m-%dT%H-%M-%SZ)
BUILD_DATE := `date -u +"%Y-%m-%d-%H-%M-%SZ"`
GOARCH    := amd64
VERSION   := 0.1.0

# Let's parse make target comments prefixed with ## and generate help output for the user. 
# Let's parse make target comments prefixed with ## and generate help output for the user. 
define PRINT_HELP_PYSCRIPT
import re, sys

for line in sys.stdin:
	match = re.match(r'^([a-zA-Z_-]+):.*?## (.*)$$', line)
	if match:
		target, help = match.groups()
		print("%-20s %s" % (target, help))
endef
export PRINT_HELP_PYSCRIPT


default: help

help:
	@python -c "$$PRINT_HELP_PYSCRIPT" < $(MAKEFILE_LIST)

# ==============================================================================
# Running tests on the local machine

test: ## Run unit tests and staticcheck locally
	go test -race -count=1 -v
	staticcheck ./...

# ==============================================================================
# Modules support

deps-reset: ## Reset go mod dependencies
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy: ## Resolve dependencies
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list: ## List Go modules
	go list -mod=mod all

build: ## Build Go binaries
	go build -ldflags "-X main.Commit=${VCS_REF} -X main.Version=${VERSION} -X main.Date=${BUILD_DATE}" -o arct ./cmd/arc/main.go

# ==============================================================================
# Go releaser

snapshot: ## GoReleaser - Build Snapshot
	goreleaser build --snapshot --rm-dist

release: ## GoReleaser - Release
	goreleaser release --rm-dist

