package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mr687/lazycopilot/pkg/config"
	"github.com/mr687/lazycopilot/pkg/copilot"
	"github.com/mr687/lazycopilot/pkg/utils"
	"github.com/spf13/cobra"
)

func newCommitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "commit",
		Short: "Generate a commit message using AI",
		Run:   commitRunner,
	}
}

func commitRunner(cmd *cobra.Command, args []string) {
	path, _ := os.Getwd()
	if !utils.IsFileExists(path) {
		fmt.Printf("Error: %s is not a valid path\n", path)
		os.Exit(1)
	}
	diff := utils.GetDiff(path)
	if diff == "" {
		os.Exit(0)
	}

	ctx := context.Background()

	commitPrompt := strings.ReplaceAll(config.COMMIT_PROMPT, "{{diff}}", diff)
	copilot := copilot.NewCopilot()
	content, err := copilot.Ask(ctx, commitPrompt, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(content)
}
