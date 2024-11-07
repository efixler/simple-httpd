BUILD_DIR := build
MODULE_NAME := simple_httpd
ifneq ($(shell command -v go >/dev/null 2>&1 && echo yes),)
    MODULE_NAME := $(shell go list -m)
endif

.DEFAULT_GOAL := build

.PHONY: fmt vet build clean help cognitive


build: vet ## build the binary, to the build/ folder (default target)
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

release-tag: latest-release-tag ## create a release tag at the next patch version. Customize with TAG_MESSAGE and/or TAG_VERSION
	$(eval TAG_VERSION ?= $(shell echo $(RELEASE_TAG) | awk -F. '{print $$1"."$$2"."$$3+1}'))
	$(eval TAG_MESSAGE ?= "Release version $(TAG_VERSION)")
	@echo "Creating release tag $(TAG_VERSION) with message: \"$(TAG_MESSAGE)\""
	@if [ "$(TAG_VERSION)" = "v0.0.0" ]; then \
        echo "Aborted. Release version cannot be 'v0.0.0'."; \
        exit 1; \
    fi
	@read -p "Continue to push this release tag? (y/n): " answer; \
    if [ "$$answer" != "y" ]; then \
        echo "Aborted."; \
        exit 1; \
    fi
	@git tag -a $(TAG_VERSION) -m '$(TAG_MESSAGE)'
	@git push origin $(TAG_VERSION)

help: ## show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

latest-release-tag: 
	$(eval RELEASE_TAG := $(shell git describe --abbrev=0 --tags $(shell git rev-list --tags --max-count=1) 2>/dev/null || echo "v0.0.0"))
	@echo "Latest release tag: $(RELEASE_TAG)"