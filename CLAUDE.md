# CLAUDE.md

This file provides guidance for AI assistants working with this codebase.

## Project Overview

This is a Go CLI application built with Cobra. It follows standard Go project layout conventions.

## Build & Test Commands

```bash
just build      # Build with version info to bin/
just install    # Install to GOBIN
just test       # Run tests
just lint       # Run golangci-lint
just check      # Run all quality checks (fmt, vet, lint, test)
just clean      # Clean build artifacts
```

## Project Structure

```
.
├── cmd/mycli/          # Main entry point (thin wrapper)
│   └── main.go         # Version injection, calls cli.Execute()
├── internal/
│   ├── cli/            # Cobra command definitions
│   │   ├── root.go     # Root command and Execute()
│   │   ├── version.go  # Version command
│   │   └── *.go        # Additional commands
│   ├── domain/         # Core business logic (no external deps)
│   ├── config/         # Configuration loading
│   └── terminal/       # Terminal output, colors, logging
├── .github/workflows/  # CI and release automation
└── docs/               # Additional documentation
```

## Key Patterns

### Adding a New Command

1. Create a new file in `internal/cli/` (e.g., `mycommand.go`)
2. Define the command with `&cobra.Command{}`
3. Add flags in `init()` using `cmd.Flags().TypeVarP()`
4. Register with `rootCmd.AddCommand(myCmd)` in `init()`

### Version Injection

Version info is injected at build time via ldflags:
- `-X main.version=...`
- `-X main.commit=...`
- `-X main.date=...`

Falls back to `debug.ReadBuildInfo()` for `go install` builds.

### Error Handling

- Commands should use `RunE` and return errors
- Wrap errors with context: `fmt.Errorf("action failed: %w", err)`
- Exit codes are handled in `cli.Execute()`

## Testing

- Tests live alongside code: `foo.go` → `foo_test.go`
- Use table-driven tests for multiple cases
- Use `t.TempDir()` for temp directories (auto-cleanup)

## Common Tasks

### Rename the CLI

1. Update `BINARY` in `justfile`
2. Update `Use` in `internal/cli/root.go`
3. Update `main` and `binary` in `.goreleaser.yaml`
4. Update module path in `go.mod` and all imports

### Add a Dependency

```bash
go get github.com/some/package
go mod tidy
```

### Create a Release

```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions runs GoReleaser automatically
```
