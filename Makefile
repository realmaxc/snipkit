BINARY  := snipkit
PKG     := ./...
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
MODULE  := github.com/realmaxc/snipkit
LDFLAGS := -s -w -X $(MODULE)/internal/version.Version=$(VERSION)

.PHONY: build test cover fmt vet lint tidy install clean help

build: ## Build the snipkit binary
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/snipkit

test: ## Run the test suite with the race detector
	go test -race $(PKG)

cover: ## Run tests and print a coverage summary
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out | tail -n 1

fmt: ## Format all Go sources
	gofmt -w .

vet: ## Run go vet
	go vet $(PKG)

lint: ## Run golangci-lint
	golangci-lint run

tidy: ## Tidy and verify go.mod/go.sum
	go mod tidy

install: ## Install snipkit into GOBIN
	go install -ldflags "$(LDFLAGS)" ./cmd/snipkit

clean: ## Remove build and coverage artifacts
	rm -f $(BINARY) coverage.out coverage.html

help: ## List available targets
	@grep -hE '^[a-z].*:.*##' $(MAKEFILE_LIST) | sort | awk -F ":.*## " '{printf "  %-10s %s\n", $$1, $$2}'
