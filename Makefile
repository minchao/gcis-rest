.PHONY: install deps lint build

GOOS := linux

install:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $$(go env GOPATH)/bin v1.12.5

deps:
	go mod download

lint:
	golangci-lint run -E gofmt ./cmd

build: ## Build the lambda binary
	GOOS=$(GOOS) go build -ldflags="-s -w" -o ./build/handler cmd/main.go
