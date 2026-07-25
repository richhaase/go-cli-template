package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/OWNER/REPO/internal/cli"
)

// Build-time variables injected via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(run())
}

// run is separate from main so deferred cleanup runs before os.Exit.
func run() int {
	// Cancel the context on Ctrl-C (SIGINT) or SIGTERM so commands
	// reading cmd.Context() can stop in-flight work cleanly.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	v, c, d := getVersionInfo()
	return cli.Execute(ctx, v, c, d)
}

// getVersionInfo returns version info, preferring build-time values
// but falling back to debug.ReadBuildInfo for `go install` builds.
func getVersionInfo() (string, string, string) {
	if version != "dev" {
		return version, commit, date
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version, commit, date
	}

	// Extract version from module info
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}

	// Extract commit, build date, and dirty status from build settings
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) >= 7 {
				commit = setting.Value[:7]
			}
		case "vcs.time":
			date = setting.Value
		case "vcs.modified":
			if setting.Value == "true" {
				commit += "-dirty"
			}
		}
	}

	return version, commit, date
}
