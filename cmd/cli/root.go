package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lazycopilot",
	Short: "lazycopilot is a tool to assist with various development tasks",
	Long:  `lazycopilot is a versatile tool designed to assist developers with a range of tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return fmt.Errorf("lazycopilot: %s is not a valid command. See 'lazycopilot --help'", args[0])
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(newAuthCommand())
	rootCmd.AddCommand(newCommitCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
