.PHONY: help build test coverage lint clean run example install-tools tidy

# Add Go bin to PATH for all targets
GOPATH ?= $(shell go env GOPATH)
export PATH := $(GOPATH)/bin:$(PATH)

# Variables
BINARY_NAME=gorest-mcp
GO=go
GOTEST=$(GO) test
GOCOVER=$(GO) tool cover
GOLINT=golangci-lint
EXAMPLE_DIR=examples/basic

# Default target
GOLANGCI_LINT_VERSION := v2.12.2

.DEFAULT_GOAL := help

help:
	@echo "GoREST-MCP Plugin v0.1.0"
	@echo ""
	@echo "Available targets:"
	@echo "  make build         - Build the plugin"
	@echo "  make test          - Run tests"
	@echo "  make coverage      - Run tests with coverage report"
	@echo "  make lint          - Run golangci-lint (bundles staticcheck, errcheck, govet, gocyclo, misspell)"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make run           - Run the basic example"
	@echo "  make example       - Build and run basic example"
	@echo "  make install-tools - Install development tools"
	@echo "  make tidy          - Tidy and verify dependencies"

# Build
build:
	@echo "Building gorest-mcp plugin..."
	$(GO) build -v ./...

# Test
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Test with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCOVER) -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Lint
lint:
	@echo "Running linter..."
	@which $(GOLINT) > /dev/null || (echo "golangci-lint not installed. Run 'make install-tools'" && exit 1)
	$(GOLINT) run --timeout=5m

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f coverage.out coverage.html
	rm -f $(EXAMPLE_DIR)/$(BINARY_NAME)
	rm -f $(EXAMPLE_DIR)/*.db
	$(GO) clean

# Run basic example
run: build
	@echo "Running basic example..."
	cd $(EXAMPLE_DIR) && $(GO) run main.go

# Build and run example
example: build
	@echo "Building example..."
	cd $(EXAMPLE_DIR) && $(GO) build -o $(BINARY_NAME) main.go
	@echo "Running example..."
	cd $(EXAMPLE_DIR) && ./$(BINARY_NAME)

# Install development tools
install-tools:
	@echo "[INFO] Installing development tools..."
	@if ! golangci-lint --version 2>/dev/null | grep -qE 'version v?2\.'; then \
		echo "  Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		GOWORK=off go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	fi
	@echo "✓ Development tools installed"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	$(GO) mod verify
	@echo "Dependencies tidied"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Full check (format, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GO) mod verify

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Show version
version:
	@echo "gorest-mcp v0.1.0"
	@$(GO) version

