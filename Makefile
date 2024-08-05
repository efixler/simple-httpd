BUILD_DIR := build
MODULE_NAME := scrape
ifneq ($(shell command -v go >/dev/null 2>&1 && echo yes),)
    MODULE_NAME := $(shell go list -m)
endif

.DEFAULT_GOAL := build

.PHONY: fmt vet build clean help cognitive test


build: vet ## build the binaries, to the build/ folder (default target)
	@echo "Building $(MODULE_NAME)..."
	@go build -o $(BUILD_DIR)/ -tags "$(TAGS)" ./...
	
clean: ## clean the build directory
	@echo "Cleaning $(MODULE_NAME)..."
	@rm -rf $(BUILD_DIR)/*

fmt: 
	@echo "Running go fmt..."
	@go fmt ./...

vet: fmt
	@echo "Running go vet and staticcheck..."
	@go vet ./...
	@staticcheck ./...

cognitive: ## run the cognitive complexity checker
	@echo "Running gocognit..."
	@gocognit  -ignore "_test|testdata" -top 5 .