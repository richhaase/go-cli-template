# Project configuration
BINARY := "mycli"
CMD_PATH := "./cmd/mycli"
BIN_DIR := "bin"

# Pinned tool versions (keep the golangci-lint version in sync with .github/workflows/ci.yml)
GOLANGCI_LINT_VERSION := "v2.12.2"

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

# Format code (rewrites files)
fmt:
    go fmt ./...
    go tool goimports -w .
    @echo "Formatted all Go files"

# Check formatting without modifying files (used by `check` / CI)
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    unformatted=$(gofmt -l .)
    if [ -n "$unformatted" ]; then
        echo "Files need formatting (run 'just fmt'):"
        echo "$unformatted"
        exit 1
    fi

# Run go vet
vet:
    go vet ./...

# Run golangci-lint (pinned version; includes staticcheck)
lint:
    go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@{{GOLANGCI_LINT_VERSION}} run

# Scan dependencies for known vulnerabilities
vuln:
    go run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Run all quality checks (non-mutating, mirrors CI)
check: fmt-check vet lint test

# Verify the GoReleaser config with a local snapshot build (no publish)
release-snapshot:
    goreleaser release --snapshot --clean

# Clean build artifacts
clean:
    rm -rf {{BIN_DIR}} dist coverage.out coverage.html
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
