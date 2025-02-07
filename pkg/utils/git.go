package utils

import (
	"os/exec"
	"strings"
)

func GetDiff(path string) string {
	if path == "" {
		path = "$(pwd)"
	}
	args := []string{
		"git",
		"-C",
		path,
		"diff",
		"--staged",
		"--no-color",
		"--no-ext-diff",
	}

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
