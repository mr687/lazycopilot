package commit

import (
	"fmt"

	"github.com/mr687/lazycopilot/pkg/config"
)

type Style string

var availableStyles map[Style]*config.CommitStyle

func init() {
	LoadStyles()
}

func LoadStyles() {
	styles := config.LoadCommitStyles()
	availableStyles = make(map[Style]*config.CommitStyle)
	for _, style := range styles {
		s := style
		availableStyles[Style(s.Name)] = &s
	}
}

func GetStylePrompt(style Style) string {
	if s, exists := availableStyles[style]; exists {
		return s.Prompt
	}
	return ""
}

func IsValidStyle(style string) bool {
	_, exists := availableStyles[Style(style)]
	return exists
}

func GetAvailableStyles() []string {
	styles := make([]string, 0, len(availableStyles))
	for style := range availableStyles {
		styles = append(styles, string(style))
	}
	return styles
}

func GetAllStyles() []config.CommitStyle {
	styles := make([]config.CommitStyle, 0, len(availableStyles))
	for _, style := range availableStyles {
		styles = append(styles, *style)
	}
	return styles
}

func AddStyle(style config.CommitStyle) error {
	styles := GetAllStyles()
	for _, s := range styles {
		if s.Name == style.Name {
			return fmt.Errorf("style '%s' already exists", style.Name)
		}
	}

	styles = append(styles, style)
	if err := config.SaveCommitStyles(styles); err != nil {
		return err
	}
	LoadStyles() // Reload styles after saving
	return nil
}

func RemoveStyle(name string) error {
	styles := GetAllStyles()
	found := false
	newStyles := make([]config.CommitStyle, 0, len(styles))

	for _, style := range styles {
		if style.Name == name {
			found = true
			continue
		}
		newStyles = append(newStyles, style)
	}

	if !found {
		return fmt.Errorf("style '%s' not found", name)
	}

	if err := config.SaveCommitStyles(newStyles); err != nil {
		return err
	}
	LoadStyles() // Reload styles after saving
	return nil
}

type StyleChange struct {
	Name   string
	Action string // "update", "add", or "keep"
}

func PreviewSyncChanges() []StyleChange {
	currentStyles := GetAllStyles()
	defaultStyleMap := make(map[string]config.CommitStyle)
	changes := make([]StyleChange, 0)

	// Create map of default styles
	for _, style := range config.DefaultCommitStyles {
		defaultStyleMap[style.Name] = style
	}

	// Check current styles
	for _, style := range currentStyles {
		if defaultStyle, exists := defaultStyleMap[style.Name]; exists {
			// Check if the content is actually different
			if style.Description != defaultStyle.Description || style.Prompt != defaultStyle.Prompt {
				changes = append(changes, StyleChange{
					Name:   style.Name,
					Action: "update",
				})
			} else {
				changes = append(changes, StyleChange{
					Name:   style.Name,
					Action: "keep",
				})
			}
		} else {
			changes = append(changes, StyleChange{
				Name:   style.Name,
				Action: "keep",
			})
		}
		delete(defaultStyleMap, style.Name)
	}

	// Check remaining default styles to be added
	for name := range defaultStyleMap {
		changes = append(changes, StyleChange{
			Name:   name,
			Action: "add",
		})
	}

	return changes
}

func SyncStyles() ([]config.CommitStyle, error) {
	currentStyles := GetAllStyles()
	defaultStyleMap := make(map[string]config.CommitStyle)

	// Create map of default styles for quick lookup
	for _, style := range config.DefaultCommitStyles {
		defaultStyleMap[style.Name] = style
	}

	// Initialize synced slice with default capacity
	synced := make([]config.CommitStyle, 0, len(currentStyles))

	// Process current styles
	for _, style := range currentStyles {
		if defaultStyle, exists := defaultStyleMap[style.Name]; exists {
			// If it's a default style, use the default version
			synced = append(synced, defaultStyle)
			delete(defaultStyleMap, style.Name)
		} else {
			// If it's not a default style, keep the custom style
			synced = append(synced, style)
		}
	}

	// Add remaining default styles that weren't in current styles
	for _, style := range defaultStyleMap {
		synced = append(synced, style)
	}

	if err := config.SaveCommitStyles(synced); err != nil {
		return nil, err
	}
	LoadStyles() // Reload styles after saving
	return synced, nil
}
