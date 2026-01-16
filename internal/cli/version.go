package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mycli %s\n", versionInfo.version)
		fmt.Printf("  commit: %s\n", versionInfo.commit)
		fmt.Printf("  built:  %s\n", versionInfo.date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
