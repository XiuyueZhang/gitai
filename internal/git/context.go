package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectContext holds information about the git project
type ProjectContext struct {
	ProjectName   string
	RecentCommits []string
	BranchName    string
	ChangedFiles  []string
	ReadmeSnippet string
	DiffStats     string
}

// GetProjectContext collects context information about the current project
func GetProjectContext() (ProjectContext, error) {
	ctx := ProjectContext{}

	// Get project name (from directory name or git remote)
	ctx.ProjectName = getProjectName()

	// Get recent commits
	commits, err := getRecentCommits(5)
	if err == nil {
		ctx.RecentCommits = commits
	}

	// Get current branch
	branch, err := getCurrentBranch()
	if err == nil {
		ctx.BranchName = branch
	}

	// Get changed files
	files, err := GetChangedFiles()
	if err == nil {
		ctx.ChangedFiles = files
	}

	// Get diff stats
	stats, err := GetDiffStats()
	if err == nil {
		ctx.DiffStats = stats
	}

	// Get README snippet
	readme := getReadmeSnippet(500)
	if readme != "" {
		ctx.ReadmeSnippet = readme
	}

	return ctx, nil
}

// getProjectName attempts to determine the project name
func getProjectName() string {
	// Try to get from git remote
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err == nil {
		url := string(output)
		// Extract repo name from URL
		parts := strings.Split(strings.TrimSpace(url), "/")
		if len(parts) > 0 {
			name := parts[len(parts)-1]
			name = strings.TrimSuffix(name, ".git")
			if name != "" {
				return name
			}
		}
	}

	// Fallback to directory name
	cwd, err := os.Getwd()
	if err == nil {
		return filepath.Base(cwd)
	}

	return "unknown"
}

// getRecentCommits returns the last N commit messages
func getRecentCommits(n int) ([]string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", n), "--pretty=format:%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 1 && commits[0] == "" {
		return []string{}, nil
	}

	return commits, nil
}

// getCurrentBranch returns the name of the current branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// getReadmeSnippet reads the first N characters from README.md if it exists
func getReadmeSnippet(maxChars int) string {
	// Try common README file names
	readmeFiles := []string{"README.md", "README.MD", "readme.md", "Readme.md"}

	for _, filename := range readmeFiles {
		content, err := os.ReadFile(filename)
		if err == nil {
			text := string(content)
			// Remove markdown headers and get first paragraph
			lines := strings.Split(text, "\n")
			var cleaned []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				cleaned = append(cleaned, line)
				if len(strings.Join(cleaned, " ")) >= maxChars {
					break
				}
			}

			snippet := strings.Join(cleaned, " ")
			if len(snippet) > maxChars {
				snippet = snippet[:maxChars] + "..."
			}
			return snippet
		}
	}

	return ""
}
