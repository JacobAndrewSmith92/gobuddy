package util

import (
	"fmt"
	"os/exec"
	"strings"
)

// checkIfGitRepo checks if the current directory is a Git repository
func checkIfGitRepo() bool {
	// Try running 'git rev-parse --is-inside-work-tree'
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// getCurrentBranch retrieves the current Git branch in the repository
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or no branch: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetBranch() (string, error) {
	if !checkIfGitRepo() {
		return "", fmt.Errorf("not a git repository")
	}
	return getCurrentBranch()
}
