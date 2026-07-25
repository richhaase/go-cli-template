package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version info set at runtime
	versionInfo struct {
		version string
		commit  string
		date    string
	}

	// Global flags
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "mycli",
		Short: "A brief description of your CLI",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
		// Execute() owns error printing and exit codes, so keep Cobra
		// from also printing errors and dumping usage on every failure.
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			// Logs go to stderr; stdout is reserved for command output.
			level := slog.LevelInfo
			if verbose {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
				Level: level,
			})))
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
)

// Execute runs the root command and returns an exit code.
// The context is canceled on SIGINT/SIGTERM (see cmd/mycli/main.go);
// commands should read cmd.Context() to honor cancellation.
func Execute(ctx context.Context, version, commit, date string) int {
	versionInfo.version = version
	versionInfo.commit = commit
	versionInfo.date = date

	// Enable the conventional `mycli --version` flag.
	rootCmd.Version = version

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose (debug) logging")

	// Local flags
	// rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")
}
