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
	CommitType       string
	Scope            string
	Diff             string
	Context          ProjectContext
	Language         string
	DetailedCommit   bool   // If true, generate multi-line commit with body
	CustomPrompt     string // Custom company/team commit guidelines
	TicketNumber     string // Ticket/issue number (e.g., JIRA-123)
	SubjectLength    string // Subject length: "short" (36 chars) or "normal" (72 chars)
	RegenerateCount  int    // Number of times regenerated (adds variation hints)
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

	// Custom company/team guidelines (if provided)
	if pb.CustomPrompt != "" {
		prompt.WriteString("COMPANY/TEAM COMMIT GUIDELINES:\n")
		prompt.WriteString(pb.CustomPrompt)
		prompt.WriteString("\n\n")
		prompt.WriteString("IMPORTANT: Follow the above guidelines strictly when generating the commit message.\n\n")
	}

	// Task description
	prompt.WriteString("TASK:\n")
	prompt.WriteString(fmt.Sprintf("Generate a %s commit message for the following changes.\n", pb.CommitType))

	if pb.Scope != "" {
		prompt.WriteString(fmt.Sprintf("Scope: %s\n", pb.Scope))
	}

	if pb.TicketNumber != "" {
		prompt.WriteString(fmt.Sprintf("Ticket/Issue Number: %s\n", pb.TicketNumber))
		prompt.WriteString(fmt.Sprintf("IMPORTANT: Include the ticket number [%s] in the commit message.\n", pb.TicketNumber))
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

	// Requirements - different based on detailed mode
	prompt.WriteString("REQUIREMENTS:\n")
	prompt.WriteString("1. Follow Conventional Commits format\n")

	// Set subject length based on configuration
	maxLength := 72
	if pb.SubjectLength == "short" {
		maxLength = 36
	}
	prompt.WriteString(fmt.Sprintf("2. Subject line: concise summary (max %d characters)\n", maxLength))

	if pb.DetailedCommit {
		// Detailed mode: include body with explanations
		prompt.WriteString("3. Body: explain WHAT changed and WHY (2-4 bullet points)\n")
		prompt.WriteString("4. Focus on the motivation and impact, not implementation details\n")
		prompt.WriteString(fmt.Sprintf("5. Use %s language\n", language))
		prompt.WriteString("6. Start subject line with lowercase letter after the type\n")
		prompt.WriteString("7. Separate subject and body with a blank line\n\n")
	} else {
		// Concise mode: subject line only
		prompt.WriteString("3. Focus on WHAT changed and WHY (concise)\n")
		prompt.WriteString(fmt.Sprintf("4. Use %s language\n", language))
		prompt.WriteString("5. Start with lowercase letter after the type\n")
		prompt.WriteString("6. Generate ONLY the subject line, no body or explanation\n\n")
	}

	// Output format
	prompt.WriteString("OUTPUT FORMAT:\n")

	// Build format string based on ticket presence
	var formatStr, exampleSubject string
	if pb.TicketNumber != "" {
		if pb.Scope != "" {
			formatStr = fmt.Sprintf("%s(%s): [%s] <subject line>", pb.CommitType, pb.Scope, pb.TicketNumber)
			exampleSubject = fmt.Sprintf("%s(%s): [%s] add user authentication endpoint", pb.CommitType, pb.Scope, pb.TicketNumber)
		} else {
			formatStr = fmt.Sprintf("%s: [%s] <subject line>", pb.CommitType, pb.TicketNumber)
			exampleSubject = fmt.Sprintf("%s: [%s] add user authentication endpoint", pb.CommitType, pb.TicketNumber)
		}
	} else {
		if pb.Scope != "" {
			formatStr = fmt.Sprintf("%s(%s): <subject line>", pb.CommitType, pb.Scope)
			exampleSubject = fmt.Sprintf("%s(%s): add user authentication endpoint", pb.CommitType, pb.Scope)
		} else {
			formatStr = fmt.Sprintf("%s: <subject line>", pb.CommitType)
			exampleSubject = fmt.Sprintf("%s: add user authentication endpoint", pb.CommitType)
		}
	}

	if pb.DetailedCommit {
		// Detailed format with body
		prompt.WriteString(formatStr + "\n\n<body with bullet points>\n\n")
		prompt.WriteString("Example:\n")
		prompt.WriteString(exampleSubject + "\n\n")
		prompt.WriteString("- Implement JWT-based authentication\n")
		prompt.WriteString("- Add login and logout endpoints\n")
		prompt.WriteString("- Include token validation middleware\n\n")
		prompt.WriteString("Generate the commit message now (subject + body with details):\n")
	} else {
		// Concise format - subject only
		prompt.WriteString(formatStr + "\n\n")
		prompt.WriteString("Example:\n")
		prompt.WriteString(exampleSubject + "\n\n")
		prompt.WriteString("Generate the commit message now (ONLY the subject line):\n")
	}

	// Add variation hint if this is a regeneration
	if pb.RegenerateCount > 0 {
		variationHints := []string{
			"Try a different perspective or emphasis in the subject line.",
			"Consider alternative wording or focus on different aspects.",
			"Rephrase with a fresh approach while maintaining accuracy.",
			"Use different verbs or structure to convey the same changes.",
			"Focus on a different aspect of the changes for variety.",
		}
		hintIndex := pb.RegenerateCount % len(variationHints)
		prompt.WriteString(fmt.Sprintf("\nNOTE: This is regeneration attempt #%d. %s\n", pb.RegenerateCount, variationHints[hintIndex]))
	}

	return prompt.String()
}

// NewPromptBuilder creates a new PromptBuilder with default values
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		Language: "en",
		Context:  ProjectContext{},
	}
}
