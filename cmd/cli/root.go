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

var rootCmd = &cobra.Command{
	Use:   "lazycopilot",
	Short: "lazycopilot is a tool to help you write better commit messages",
	Long:  `lazycopilot is a tool to help you write better commit messages`,
	Run:   runner,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&config.Path, "path", "p", "", "Path to the git repository")
}

func runner(cmd *cobra.Command, args []string) {
	if config.Path == "" {
		config.Path, _ = os.Getwd()
	}
	if !utils.IsFileExists(config.Path) {
		fmt.Printf("Error: %s is not a valid path\n", config.Path)
		os.Exit(1)
	}
	diff := utils.GetDiff(config.Path)
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
