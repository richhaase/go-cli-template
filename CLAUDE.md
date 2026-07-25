# CLAUDE.md

This file provides guidance for AI assistants working with this codebase.

## Project Overview

This is a Go CLI application built with Cobra. It follows standard Go project layout conventions.

## Build & Test Commands

```bash
just build             # Build with version info to bin/
just install           # Install to GOBIN
just test              # Run tests
just lint              # Run golangci-lint (pinned version)
just check             # Run all quality checks (fmt-check, vet, lint, test) — non-mutating
just fmt               # Format code (rewrites files)
just vuln              # Run govulncheck
just release-snapshot  # Verify GoReleaser config locally (no publish)
just deps-update       # Update all dependencies to latest
just clean             # Clean build artifacts
```

## Project Structure

```
.
├── cmd/mycli/          # Main entry point (thin wrapper)
│   └── main.go         # Version injection, signal handling, calls cli.Execute()
├── internal/
│   ├── cli/            # Cobra command definitions
│   │   ├── root.go     # Root command, Execute(), slog setup
│   │   ├── version.go  # Version command
│   │   └── *.go        # Additional commands
│   ├── domain/         # Core business logic (no external deps)
│   ├── config/         # Configuration loading
│   └── terminal/       # TTY detection helpers
├── scripts/            # Bootstrap tooling (rename.sh)
├── .github/workflows/  # CI and release automation
└── docs/               # Additional documentation
```

## Key Patterns

### Adding a New Command

1. Create a new file in `internal/cli/` (e.g., `mycommand.go`)
2. Define the command with `&cobra.Command{}`
3. Add flags in `init()` using `cmd.Flags().TypeVarP()`
4. Register with `rootCmd.AddCommand(myCmd)` in `init()`

### Context and Cancellation

- `main()` installs `signal.NotifyContext` (SIGINT/SIGTERM) and calls
  `rootCmd.ExecuteContext(ctx)`
- Long-running commands should read `cmd.Context()` and stop when it is
  canceled (see `internal/cli/example.go`)

### Logging

- Use `log/slog`; the default handler writes to stderr and is configured in
  the root command's `PersistentPreRun`
- The persistent `--verbose`/`-v` flag switches the level from Info to Debug
- stdout is reserved for command output; logs always go to stderr

### Output

- Write command output with `fmt.Fprintf(cmd.OutOrStdout(), ...)` rather than
  `fmt.Printf` so tests can capture it

### Version Injection

Version info is injected at build time via ldflags:
- `-X main.version=...`
- `-X main.commit=...`
- `-X main.date=...`

Falls back to `debug.ReadBuildInfo()` for `go install` builds.
`--version` is also supported (rootCmd.Version is set in Execute()).

### Error Handling

- Commands should use `RunE` and return errors
- Wrap errors with context: `fmt.Errorf("action failed: %w", err)`
- Exit codes are handled in `cli.Execute()`; the root command sets
  `SilenceErrors`/`SilenceUsage` so errors are printed exactly once

## Testing

- Tests live alongside code: `foo.go` → `foo_test.go`
- Use table-driven tests for multiple cases (see `internal/cli/example_test.go`)
- Drive Cobra commands via `rootCmd.SetArgs` + `SetOut`/`SetErr` buffers
- Use `t.TempDir()` for temp directories and `t.Setenv()` for env vars
  (see `internal/config/config_test.go`)

## Common Tasks

### Rename the CLI

Use the bootstrap script, which handles all of the below (module path,
binary name, env-var prefix, README cleanup, LICENSE):

```bash
./scripts/rename.sh -o myuser -r my-cli
```

Manual equivalent:
1. Update `BINARY` in `justfile`
2. Update `Use` in `internal/cli/root.go`
3. Update `main` and `binary` in `.goreleaser.yaml`
4. Update module path in `go.mod` and all imports
5. Update the `MYCLI_*` env-var prefix in `internal/config/config.go`

### Add a Dependency

```bash
go get github.com/some/package
go mod tidy
```

### Create a Release

```bash
just release-snapshot   # local dry-run of the GoReleaser config
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions runs GoReleaser automatically
```
