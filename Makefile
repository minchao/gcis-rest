.PHONY: help install deps lint build clean package deploy

PROJECT_NAME ?= gcis-rest
CFN_BUCKET_NAME ?= $(PROJECT_NAME)-bucket
CFN_TEMPLATE := ./template.yml
CFN_PACKAGED_TEMPLATE := ./build/packaged.yml
CFN_BUILD_DIR := $(shell dirname $(CFN_PACKAGED_TEMPLATE))
CFN_PARAMETER_FILE ?=
CFN_PARAMETER_OVERRIDES := $(if $(CFN_PARAMETER_FILE:""=),--parameter-overrides $(shell jq -j '.[] | "\"" + .ParameterKey + "=" + .ParameterValue +"\" "' $(CFN_PARAMETER_FILE)),)
GOOS := linux

help:
	@echo "\nUsage: $ make COMMAND\n\nCommands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'

install: ## Install required tools
	pip install --user --upgrade awscli aws-sam-cli
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $$(go env GOPATH)/bin v1.12.5

deps: ## Install dependencies
	go mod download

lint: ## Run all go linters
	golangci-lint run -E gofmt ./cmd/... ./internal/...

build: ## Build the lambda binary
	GOOS=$(GOOS) go build -ldflags="-s -w" -o ./build/company cmd/company/main.go

clean: ## Clean all artifacts
	rm -rf ./build

package:build ## Package SAM template
	sam package \
		--template-file $(CFN_TEMPLATE) \
		--output-template-file $(CFN_PACKAGED_TEMPLATE) \
		--s3-bucket $(CFN_BUCKET_NAME)

deploy:package ## Deploy packaged SAM template
	sam deploy \
		--template-file $(CFN_PACKAGED_TEMPLATE) \
		--stack-name $(PROJECT_NAME) \
		--capabilities CAPABILITY_IAM \
		$(CFN_PARAMETER_OVERRIDES)
