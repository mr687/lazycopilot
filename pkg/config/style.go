package config

import (
	"os"
	"path/filepath"

	"github.com/mr687/lazycopilot/pkg/utils"
)

type CommitStyle struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

var DefaultCommitStyles = []CommitStyle{
	{
		Name:        "normal",
		Description: "Standard commit message style",
		Prompt:      "\n\nWriting style clear, concise, and to the point. Focus on the technical change and why it was made.",
	},
	{
		Name:        "funny",
		Description: "Funny commit messages",
		Prompt:      "\n\nWriting style humorous, lighthearted, and potentially self-deprecating. Focus on making the reader smile while still conveying the essence of the change.",
	},
	{
		Name:        "wise",
		Description: "Wise and inspirational commit messages",
		Prompt:      "\n\nWriting style wise, inspirational, and potentially poetic. Focus on providing a deeper meaning to the technical change.",
	},
	{
		Name:        "trolling",
		Description: "Trolling commit messages",
		Prompt:      "\n\nWriting style playful, slightly provocative, and potentially sarcastic. Focus on A bit of a jab, but still conveys the technical change.",
	},
}

func GetCommitStylesConfigPath() string {
	configDir := utils.GetConfigPath()
	if configDir == "" {
		home := os.Getenv("HOME")
		if home == "" {
			return ""
		}
		configDir = filepath.Join(home, DEFAULT_APP_PATHS)
	}
	return filepath.Join(configDir, APP_DIR_NAME, STYLES_FILE_NAME)
}

func LoadCommitStyles() []CommitStyle {
	configPath := GetCommitStylesConfigPath()
	if configPath == "" {
		return DefaultCommitStyles
	}

	var styles []CommitStyle
	err := utils.LoadFileJson(configPath, &styles)
	if err != nil {
		// If config doesn't exist, create it with defaults
		_ = utils.SaveFile(configPath, DefaultCommitStyles)
		return DefaultCommitStyles
	}

	return styles
}

func SaveCommitStyles(styles []CommitStyle) error {
	configPath := GetCommitStylesConfigPath()
	if configPath == "" {
		return nil
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return utils.SaveFile(configPath, styles)
}
