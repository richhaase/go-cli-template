// Package terminal provides utilities for terminal output,
// including color support and TTY detection.
package terminal

import (
	"os"

	"golang.org/x/term"
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// IsStdoutTTY returns true if stdout is a terminal.
func IsStdoutTTY() bool {
	return IsTerminal(int(os.Stdout.Fd()))
}

// IsStderrTTY returns true if stderr is a terminal.
func IsStderrTTY() bool {
	return IsTerminal(int(os.Stderr.Fd()))
}
