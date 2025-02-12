package cli

import (
	"fmt"

	"github.com/mr687/lazycopilot/pkg/commit"
	"github.com/mr687/lazycopilot/pkg/config"
	"github.com/spf13/cobra"
)

func newStyleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "style",
		Short: "Manage commit message styles",
	}

	cmd.AddCommand(
		newStyleListCommand(),
		newStyleAddCommand(),
		newStyleRemoveCommand(),
	)

	return cmd
}

func newStyleListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available commit styles",
		Run: func(cmd *cobra.Command, args []string) {
			styles := commit.GetAllStyles()
			for _, style := range styles {
				fmt.Printf("%s: %s\n", style.Name, style.Description)
			}
		},
	}
}

func newStyleAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <description> <prompt>",
		Short: "Add a new commit style",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			description := args[1]
			prompt := args[2]

			style := config.CommitStyle{
				Name:        name,
				Description: description,
				Prompt:      prompt,
			}

			if err := commit.AddStyle(style); err != nil {
				fmt.Printf("Error adding style: %v\n", err)
				return
			}
			fmt.Printf("Successfully added style '%s'\n", name)
		},
	}
	return cmd
}

func newStyleRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a commit style",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if err := commit.RemoveStyle(name); err != nil {
				fmt.Printf("Error removing style: %v\n", err)
				return
			}
			fmt.Printf("Successfully removed style '%s'\n", name)
		},
	}
}
