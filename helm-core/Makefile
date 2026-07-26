.PHONY: build build-linux build-all run vet clean help

# Binary name
BINARY=helm

# Build directories
BUILD_DIR=./bin
DIST_DIR=./dist
CMD_DIR=./cmd/helm

# Go parameters
GOFLAGS=-trimpath
LDFLAGS=-s -w

# Version info
VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

## build: Build the binary for the current platform (requires Linux)
build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY)"

## build-linux: Cross-compile for Linux amd64
build-linux:
	@echo "Cross-compiling $(BINARY) for linux/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 $(CMD_DIR)
	@echo "Built: $(BUILD_DIR)/$(BINARY)-linux-amd64"

## build-all: Cross-compile for all supported platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY)-linux-arm64 $(CMD_DIR)
	@echo "Built all targets in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

## run: Build and run the server
run: build
	@echo "Starting $(BINARY)..."
	$(BUILD_DIR)/$(BINARY)

## vet: Run static analysis
vet:
	go vet ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Done."

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
