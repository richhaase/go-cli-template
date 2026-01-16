# Project configuration
BINARY := "mycli"
CMD_PATH := "./cmd/mycli"
BIN_DIR := "bin"

# Default recipe: show available commands
default:
    @just --list

# Build the binary with version info
build:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    mkdir -p {{BIN_DIR}}
    go build -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
        -o {{BIN_DIR}}/{{BINARY}} {{CMD_PATH}}
    echo "Built {{BIN_DIR}}/{{BINARY}} (${VERSION})"

# Install to GOBIN with version info
install:
    #!/usr/bin/env bash
    set -euo pipefail
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
    DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    go install -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" {{CMD_PATH}}
    echo "Installed {{BINARY}} to $(go env GOBIN || echo '$(go env GOPATH)/bin')"

# Run the CLI (builds first if needed)
run *ARGS:
    @just build
    ./{{BIN_DIR}}/{{BINARY}} {{ARGS}}

# Run all tests
test:
    go test -race ./...

# Run tests with coverage
test-coverage:
    go test -race -covermode=atomic -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report: coverage.html"

# Format code
fmt:
    go fmt ./...
    @echo "Formatted all Go files"

# Run go vet
vet:
    go vet ./...

# Run golangci-lint
lint:
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run

# Run staticcheck
staticcheck:
    go run honnef.co/go/tools/cmd/staticcheck@latest ./...

# Run all quality checks
check: fmt vet lint test

# Clean build artifacts
clean:
    rm -rf {{BIN_DIR}} coverage.out coverage.html
    @echo "Cleaned build artifacts"

# Update dependencies
deps:
    go mod tidy
    go mod verify

# Generate (if you have go:generate directives)
generate:
    go generate ./...

# Show module dependencies
deps-list:
    go list -m all

# Update all dependencies to latest
deps-update:
    go get -u ./...
    go mod tidy
