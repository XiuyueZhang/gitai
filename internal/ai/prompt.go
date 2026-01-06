package ai

import (
	"fmt"
	"strings"
)

// ProjectContext holds information about the project for context-aware commit messages
type ProjectContext struct {
	ProjectName   string
	RecentCommits []string
	BranchName    string
	ChangedFiles  []string
	ReadmeSnippet string
	DiffStats     string
}

// PromptBuilder constructs prompts for AI commit message generation
type PromptBuilder struct {
	CommitType string
	Scope      string
	Diff       string
	Context    ProjectContext
	Language   string
}

// Build constructs the complete prompt for Ollama
func (pb *PromptBuilder) Build() string {
	var prompt strings.Builder

	// Header
	prompt.WriteString("You are a Git commit message generator expert.\n\n")

	// Project context
	if pb.Context.ProjectName != "" || pb.Context.BranchName != "" || len(pb.Context.RecentCommits) > 0 {
		prompt.WriteString("PROJECT CONTEXT:\n")

		if pb.Context.ProjectName != "" {
			prompt.WriteString(fmt.Sprintf("- Project: %s\n", pb.Context.ProjectName))
		}

		if pb.Context.BranchName != "" {
			prompt.WriteString(fmt.Sprintf("- Branch: %s\n", pb.Context.BranchName))
		}

		if len(pb.Context.RecentCommits) > 0 {
			prompt.WriteString("- Recent commits style:\n")
			for _, commit := range pb.Context.RecentCommits {
				prompt.WriteString(fmt.Sprintf("  * %s\n", commit))
			}
		}

		if pb.Context.ReadmeSnippet != "" {
			prompt.WriteString(fmt.Sprintf("- Project description: %s\n", pb.Context.ReadmeSnippet))
		}

		prompt.WriteString("\n")
	}

	// Task description
	prompt.WriteString("TASK:\n")
	prompt.WriteString(fmt.Sprintf("Generate a %s commit message for the following changes.\n", pb.CommitType))

	if pb.Scope != "" {
		prompt.WriteString(fmt.Sprintf("Scope: %s\n", pb.Scope))
	}

	language := pb.Language
	if language == "" {
		language = "en"
	}
	prompt.WriteString(fmt.Sprintf("Language: %s\n\n", language))

	// Changed files
	if len(pb.Context.ChangedFiles) > 0 {
		prompt.WriteString("CHANGED FILES:\n")
		for _, file := range pb.Context.ChangedFiles {
			prompt.WriteString(fmt.Sprintf("- %s\n", file))
		}
		prompt.WriteString("\n")
	}

	// Diff stats
	if pb.Context.DiffStats != "" {
		prompt.WriteString("CHANGES SUMMARY:\n")
		prompt.WriteString(pb.Context.DiffStats)
		prompt.WriteString("\n\n")
	}

	// Actual diff (truncated if too long)
	prompt.WriteString("CHANGES:\n")
	diff := pb.Diff
	// Limit diff to ~2000 characters to avoid token limits
	if len(diff) > 2000 {
		diff = diff[:2000] + "\n... (truncated)"
	}
	prompt.WriteString(diff)
	prompt.WriteString("\n\n")

	// Requirements
	prompt.WriteString("REQUIREMENTS:\n")
	prompt.WriteString("1. Follow Conventional Commits format\n")
	prompt.WriteString("2. Be concise (max 72 characters for subject line)\n")
	prompt.WriteString("3. Focus on WHAT changed and WHY (not HOW)\n")
	prompt.WriteString(fmt.Sprintf("4. Use %s language\n", language))
	prompt.WriteString("5. Start with lowercase letter after the type\n")
	prompt.WriteString("6. Do not include any explanation, just the commit message\n\n")

	// Output format
	prompt.WriteString("OUTPUT FORMAT:\n")
	if pb.Scope != "" {
		prompt.WriteString(fmt.Sprintf("%s(%s): <message>\n\n", pb.CommitType, pb.Scope))
		prompt.WriteString("Example:\n")
		prompt.WriteString(fmt.Sprintf("%s(%s): add user authentication endpoint\n\n", pb.CommitType, pb.Scope))
	} else {
		prompt.WriteString(fmt.Sprintf("%s: <message>\n\n", pb.CommitType))
		prompt.WriteString("Example:\n")
		prompt.WriteString(fmt.Sprintf("%s: add user authentication endpoint\n\n", pb.CommitType))
	}

	prompt.WriteString("Generate the commit message now (ONLY the commit message, no explanation):\n")

	return prompt.String()
}

// NewPromptBuilder creates a new PromptBuilder with default values
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		Language: "en",
		Context:  ProjectContext{},
	}
}
