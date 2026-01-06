package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	diff := string(output)
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found\nStage your changes first:\n  $ git add <files>")
	}

	return diff, nil
}

// GetChangedFiles returns list of files with staged changes
func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}

// GetDiffStats returns statistics about the changes
func GetDiffStats() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--stat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff stats: %w", err)
	}

	return string(output), nil
}

// GetChangedFilesWithStats returns files with their change statistics
func GetChangedFilesWithStats() ([]FileChange, error) {
	cmd := exec.Command("git", "diff", "--cached", "--numstat")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	var changes []FileChange
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		changes = append(changes, FileChange{
			File:      parts[2],
			Additions: parts[0],
			Deletions: parts[1],
		})
	}

	return changes, nil
}

// FileChange represents statistics for a single file
type FileChange struct {
	File      string
	Additions string
	Deletions string
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}
