package main

import (
	"os"
	"runtime/debug"

	"github.com/OWNER/REPO/internal/cli"
)

// Build-time variables injected via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	v, c, d := getVersionInfo()
	os.Exit(cli.Execute(v, c, d))
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

	// Extract commit and dirty status from build settings
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) >= 7 {
				commit = setting.Value[:7]
			}
		case "vcs.modified":
			if setting.Value == "true" {
				commit += "-dirty"
			}
		}
	}

	return version, commit, date
}
