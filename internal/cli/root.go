package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	versionInfo struct {
		version string
		commit  string
		date    string
	}

	verbose bool

	rootCmd = &cobra.Command{
		Use:   "mycli",
		Short: "A brief description of your CLI",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			level := slog.LevelInfo
			if verbose {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(cmd.ErrOrStderr(), &slog.HandlerOptions{
				Level: level,
			})))
		},
	}
)

// Execute runs the root command and returns an exit code.
func Execute(ctx context.Context, version, commit, date string) int {
	versionInfo.version = version
	versionInfo.commit = commit
	versionInfo.date = date

	rootCmd.Version = version

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose (debug) logging")

}
