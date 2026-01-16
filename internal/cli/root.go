package cli

import (
	"fmt"
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

	rootCmd = &cobra.Command{
		Use:   "mycli",
		Short: "A brief description of your CLI",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
)

// Execute runs the root command and returns an exit code.
func Execute(version, commit, date string) int {
	versionInfo.version = version
	versionInfo.commit = commit
	versionInfo.date = date

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}
	return 0
}

func init() {
	// Global flags
	// rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Local flags
	// rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "config file path")
}
