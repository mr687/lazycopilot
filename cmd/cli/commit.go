package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mr687/lazycopilot/pkg/commit"
	"github.com/mr687/lazycopilot/pkg/config"
	"github.com/mr687/lazycopilot/pkg/copilot"
	"github.com/mr687/lazycopilot/pkg/utils"
	"github.com/spf13/cobra"
)

type TitleStyle int

const (
	Normal TitleStyle = iota
	Funny
	Wise
	Trolling
)

func newCommitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Git commit related commands",
	}

	cmd.AddCommand(
		newCommitGenCommand(),
		newCommitStyleListCommand(),
		newCommitStyleAddCommand(),
		newCommitStyleRemoveCommand(),
		newCommitStyleSyncCommand(),
	)

	return cmd
}

func newCommitGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen [repository-path]",
		Short: "Generate a commit message using AI",
		Run:   commitRunner,
	}
	cmd.Flags().StringP("path", "p", "", "Path to the repository (default is current directory)")
	cmd.Flags().BoolP("stage", "s", false, "Stage changes if no staged changes are detected")
	cmd.Flags().BoolP("title-only", "t", false, "Generate only the commit title")
	cmd.Flags().StringP("style", "S", "normal", fmt.Sprintf("Style of the commit title: %s", strings.Join(commit.GetAvailableStyles(), ", ")))
	cmd.Flags().BoolP("no-commit", "n", false, "Do not commit the generated content immediately")
	return cmd
}

func newCommitStyleListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "styles",
		Short: "List all available commit styles",
		Run: func(cmd *cobra.Command, args []string) {
			styles := commit.GetAllStyles()
			for _, style := range styles {
				fmt.Printf("%s: %s\n", style.Name, style.Description)
			}
		},
	}
}

func newCommitStyleAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "style-add <name> <description> <prompt>",
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
}

func newCommitStyleRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "style-remove <name>",
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

func newCommitStyleSyncCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "style-sync",
		Short: "Sync commit styles with default styles",
		Run: func(cmd *cobra.Command, args []string) {
			changes := commit.PreviewSyncChanges()

			// Count actual changes (excluding "keep" actions)
			changeCount := 0
			for _, change := range changes {
				if change.Action != "keep" {
					changeCount++
				}
			}

			if changeCount == 0 {
				fmt.Println("No changes needed. All styles are up to date.")
				return
			}

			// Show preview
			fmt.Println("The following changes will be made:")
			for _, change := range changes {
				switch change.Action {
				case "update":
					fmt.Printf("  • Update style: %s\n", change.Name)
				case "add":
					fmt.Printf("  • Add style: %s\n", change.Name)
				}
			}

			// Ask for confirmation
			fmt.Print("\nDo you want to proceed? [y/N] ")
			var response string
			fmt.Scanln(&response)

			if strings.ToLower(response) != "y" {
				fmt.Println("Sync cancelled.")
				return
			}

			// Proceed with sync
			styles, err := commit.SyncStyles()
			if err != nil {
				fmt.Printf("Error syncing styles: %v\n", err)
				return
			}
			fmt.Printf("Successfully synced %d commit styles\n", len(styles))
		},
	}
}

func commitRunner(cmd *cobra.Command, args []string) {
	path, _ := cmd.Flags().GetString("path")
	if path == "" {
		path, _ = os.Getwd()
	}
	if !utils.IsFileExists(path) {
		fmt.Printf("Error: The specified path '%s' is not valid or does not exist.\n", path)
		os.Exit(1)
	}

	var diff string

	stage, _ := cmd.Flags().GetBool("stage")
	if stage {
		utils.StageChanges(path)
		diff = utils.GetDiff(path, true)
		if diff == "" {
			fmt.Println("Error: No changes detected to commit after staging. Please make sure you have changes to commit.")
			os.Exit(1)
		}
	}

	if diff == "" {
		diff = utils.GetDiff(path, true)
		if diff == "" {
			fmt.Println("Error: No staged changes detected. Use the --stage flag to stage all changes before committing.")
			os.Exit(1)
		}
	}

	ctx := context.Background()

	style, _ := cmd.Flags().GetString("style")
	if !commit.IsValidStyle(style) {
		fmt.Printf("Error: Invalid style '%s'. Available styles: %s\n", style, strings.Join(commit.GetAvailableStyles(), ", "))
		os.Exit(1)
	}

	commitPrompt := strings.ReplaceAll(config.COMMIT_PROMPT, "{{diff}}", diff)
	titleOnly, _ := cmd.Flags().GetBool("title-only")
	if titleOnly {
		commitPrompt += "\n\nGenerate only the commit title."
	}

	if stylePrompt := commit.GetStylePrompt(commit.Style(style)); stylePrompt != "" {
		commitPrompt += stylePrompt
	}

	copilot := copilot.NewCopilot()
	content, err := copilot.Ask(ctx, commitPrompt, nil)
	if err != nil {
		fmt.Printf("Error: Failed to generate commit message. Details: %v\n", err)
		os.Exit(1)
	}

	// Remove code block wrappers from the generated content
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	// Remove extra newlines from the generated content
	content = strings.Trim(content, "\n")

	noCommit, _ := cmd.Flags().GetBool("no-commit")
	if !noCommit {
		commitMessage := strings.SplitN(content, "\n", 2)
		commitTitle := commitMessage[0]
		var commitBody string
		if len(commitMessage) > 1 {
			commitBody = commitMessage[1]
		}

		commitFile, err := os.CreateTemp("", "commitmsg")
		if err != nil {
			fmt.Printf("Error: Failed to create temporary file for commit message. Details: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(commitFile.Name())

		commitFileContent := commitTitle
		if commitBody != "" {
			commitFileContent += "\n\n" + commitBody
		}
		if _, err := commitFile.WriteString(commitFileContent); err != nil {
			fmt.Printf("Error: Failed to write to temporary file for commit message. Details: %v\n", err)
			os.Exit(1)
		}
		commitFile.Close()

		commitArgs := []string{"git", "-C", path, "commit", "-e", "-F", commitFile.Name()}
		commitCmd := exec.Command(commitArgs[0], commitArgs[1:]...)
		commitCmd.Stdin = os.Stdin
		commitCmd.Stdout = os.Stdout
		commitCmd.Stderr = os.Stderr
		err = commitCmd.Run()
		if err != nil {
			fmt.Printf("Error: Failed to commit changes. Details: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(content)
	}
}
