package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Flags for the example command
	exampleName string
	exampleCount int
)

var exampleCmd = &cobra.Command{
	Use:   "example [message]",
	Short: "An example command to demonstrate patterns",
	Long: `This is an example command that demonstrates common CLI patterns:
- Positional arguments
- Flag parsing
- Error handling`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := "Hello"
		if len(args) > 0 {
			message = args[0]
		}

		for i := 0; i < exampleCount; i++ {
			fmt.Printf("%s, %s!\n", message, exampleName)
		}

		return nil
	},
}

func init() {
	exampleCmd.Flags().StringVarP(&exampleName, "name", "n", "World", "name to greet")
	exampleCmd.Flags().IntVarP(&exampleCount, "count", "c", 1, "number of times to repeat")

	rootCmd.AddCommand(exampleCmd)
}
