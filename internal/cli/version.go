package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "mycli %s\n", versionInfo.version)
		fmt.Fprintf(out, "  commit: %s\n", versionInfo.commit)
		fmt.Fprintf(out, "  built:  %s\n", versionInfo.date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
