package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		Use:   "commit [repository-path]",
		Short: "Generate a commit message using AI",
		Run:   commitRunner,
	}
	cmd.Flags().StringP("path", "p", "", "Path to the repository (default is current directory)")
	cmd.Flags().BoolP("stage", "s", false, "Stage changes if no staged changes are detected")
	cmd.Flags().BoolP("title-only", "t", false, "Generate only the commit title")
	cmd.Flags().StringP("style", "S", "normal", "Style of the commit title: normal, funny, wise, trolling")
	cmd.Flags().BoolP("no-commit", "n", false, "Do not commit the generated content immediately")
	return cmd
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

	commitPrompt := strings.ReplaceAll(config.COMMIT_PROMPT, "{{diff}}", diff)
	titleOnly, _ := cmd.Flags().GetBool("title-only")
	if titleOnly {
		commitPrompt += "\n\nGenerate only the commit title."
	}
	style, _ := cmd.Flags().GetString("style")
	switch style {
	case "funny":
		commitPrompt += "\n\nMake the commit title funny and add emojis."
	case "wise":
		commitPrompt += "\n\nMake the commit title wise and inspirational."
	case "trolling":
		commitPrompt += "\n\nMake the commit title trolling and sarcastic."
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
