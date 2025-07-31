# PTP Configuration Parser Makefile

# Variables
BINARY_NAME=ptp-config-parser
VERSION?=1.0.0
BUILD_DIR=bin
MAIN_FILES=main.go types.go
GO_VERSION=1.21

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"
BUILD_FLAGS=-v $(LDFLAGS)

# Default target
.DEFAULT_GOAL := help

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: build
build: clean ## Build the application
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILES)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-linux
build-linux: clean ## Build for Linux
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILES)
	@echo "✅ Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

.PHONY: build-windows
build-windows: clean ## Build for Windows
	@echo "Building $(BINARY_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILES)
	@echo "✅ Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

.PHONY: build-darwin
build-darwin: clean ## Build for macOS
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILES)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILES)
	@echo "✅ macOS builds complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-*"

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms
	@echo "✅ All platform builds complete"

##@ Testing

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests complete"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -count=1 ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

##@ Code Quality

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✅ Code formatted"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "✅ Vet complete"

.PHONY: lint
lint: ## Run golint (requires golint to be installed)
	@echo "Running golint..."
	@command -v golint >/dev/null 2>&1 || { echo >&2 "golint not installed. Run: go install golang.org/x/lint/golint@latest"; exit 1; }
	golint ./...
	@echo "✅ Lint complete"

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)
	@echo "✅ All checks passed"

##@ Dependencies

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

##@ Running

.PHONY: run
run: ## Run the application with bidirectional example
	@echo "Running $(BINARY_NAME) with bidirectional example..."
	$(GOCMD) run $(MAIN_FILES) examples/bidirectional.yaml

.PHONY: run-example
run-example: ## Run with a specific example (usage: make run-example EXAMPLE=filename.yaml)
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-example EXAMPLE=filename.yaml"; \
		echo "Available examples:"; \
		ls examples/*.yaml 2>/dev/null || echo "No examples found"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) with $(EXAMPLE)..."
	$(GOCMD) run $(MAIN_FILES) examples/$(EXAMPLE)

.PHONY: install
install: build ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

##@ Cleanup

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

.PHONY: clean-deps
clean-deps: ## Clean dependency cache
	@echo "Cleaning dependency cache..."
	$(GOCMD) clean -modcache
	@echo "✅ Dependency cache cleaned"

##@ Docker (Optional)

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	docker tag $(BINARY_NAME):$(VERSION) $(BINARY_NAME):latest
	@echo "✅ Docker image built: $(BINARY_NAME):$(VERSION)"

.PHONY: docker-run
docker-run: ## Run application in Docker
	@echo "Running $(BINARY_NAME) in Docker..."
	docker run --rm -v $(PWD)/examples:/app/examples $(BINARY_NAME):latest examples/bidirectional.yaml

##@ Information

.PHONY: info
info: ## Display project information
	@echo "Project Information:"
	@echo "  Name: $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Build Directory: $(BUILD_DIR)"
	@echo "  Main Files: $(MAIN_FILES)"

.PHONY: version
version: ## Display version
	@echo "$(BINARY_NAME) v$(VERSION)"

# Development workflow helpers
.PHONY: dev
dev: clean fmt check build ## Complete development workflow (clean, format, check, build)
	@echo "✅ Development workflow complete"

.PHONY: release
release: clean check test build-all ## Release workflow (clean, check, test, build all platforms)
	@echo "✅ Release workflow complete"
	@echo "Release artifacts:"
	@ls -la $(BUILD_DIR)/

# Validation helpers
.PHONY: validate-examples
validate-examples: build ## Validate all example files
	@echo "Validating example files..."
	@for file in examples/*.yaml; do \
		if [ -f "$$file" ] && [ -s "$$file" ]; then \
			echo "Validating $$file..."; \
			./$(BUILD_DIR)/$(BINARY_NAME) "$$file" || exit 1; \
		fi; \
	done
	@echo "✅ All examples validated"

# Quick development commands
.PHONY: quick-test
quick-test: fmt test ## Quick test (format and test only)

.PHONY: quick-build
quick-build: ## Quick build without cleaning
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILES)