package utils

import (
	"os/exec"
	"strings"
)

func GetDiff(path string, staged bool) string {
	if path == "" {
		path = "$(pwd)"
	}
	args := []string{
		"git",
		"-C",
		path,
		"diff",
	}
	if staged {
		args = append(args, "--staged")
	}
	args = append(args, "--no-color", "--no-ext-diff")

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func StageChanges(path string) error {
	if path == "" {
		path = "$(pwd)"
	}
	args := []string{
		"git",
		"-C",
		path,
		"add",
		".",
	}

	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	return err
}
