# Architecture

## Overview

This CLI follows standard Go project layout with the following structure:

```
├── cmd/mycli/      - Application entry point
├── internal/       - Private application code
│   ├── cli/        - Cobra commands
│   ├── config/     - Configuration management
│   ├── domain/     - Core business logic
│   └── terminal/   - Terminal utilities
└── docs/           - Documentation
```

## Design Principles

### Separation of Concerns

- **cmd/**: Thin entry point, handles version injection
- **cli/**: Command definitions and flag parsing
- **domain/**: Pure business logic, no external dependencies
- **config/**: Configuration loading from multiple sources
- **terminal/**: Terminal I/O, colors, formatting

### Dependency Direction

```
cmd → cli → domain
         → config
         → terminal
```

The `domain` package should have zero external dependencies and contain only pure Go code.

## Adding Features

### New Command

1. Create `internal/cli/commandname.go`
2. Define command struct and `init()` to register it
3. Call domain logic from command handler

### New Domain Logic

1. Add types to `internal/domain/types.go`
2. Add logic to appropriate file in `internal/domain/`
3. Write tests alongside implementation

## Configuration

Configuration follows three-tier precedence:
1. **Flags** - Command-line arguments (highest priority)
2. **Environment** - `MYCLI_*` environment variables
3. **Config file** - `.mycli.yaml` or `~/.config/mycli/config.yaml`
4. **Defaults** - Built-in defaults (lowest priority)
