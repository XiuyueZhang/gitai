package git

import (
	"fmt"
	"os/exec"
)

// CommitWithMessage creates a git commit with the given message
func CommitWithMessage(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// GetLastCommit returns the last commit hash and message
func GetLastCommit() (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%H %s")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get last commit: %w", err)
	}

	return string(output), nil
}
