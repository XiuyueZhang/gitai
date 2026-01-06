package git

import (
	"regexp"
	"strings"
)

// ExtractTicketFromBranch attempts to extract ticket number from branch name
// Supports patterns like:
// - feature/PROJ-123-description
// - bugfix/JIRA-456-fix-bug
// - PROJ-789
func ExtractTicketFromBranch(branchName, pattern string) string {
	if branchName == "" {
		return ""
	}

	// If custom pattern provided, use it
	if pattern != "" {
		re, err := regexp.Compile(pattern)
		if err == nil {
			matches := re.FindStringSubmatch(branchName)
			if len(matches) > 0 {
				return matches[0]
			}
		}
	}

	// Default patterns for common ticket formats
	patterns := []string{
		`[A-Z]+-\d+`,           // PROJ-123, JIRA-456
		`[A-Z]{2,}-\d+`,        // ABC-123, ABCD-456
		`#\d+`,                 // GitHub/GitLab issues: #123
		`GH-\d+`,               // GitHub: GH-123
		`[A-Z]+_\d+`,           // PROJ_123
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		matches := re.FindStringSubmatch(branchName)
		if len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// FormatTicketNumber formats ticket number with prefix if needed
func FormatTicketNumber(ticket, prefix string) string {
	ticket = strings.TrimSpace(ticket)
	if ticket == "" {
		return ""
	}

	// If ticket already has prefix, return as is
	if strings.Contains(ticket, "-") || strings.Contains(ticket, "#") {
		return ticket
	}

	// If only number provided and prefix exists, combine them
	if prefix != "" && regexp.MustCompile(`^\d+$`).MatchString(ticket) {
		return prefix + "-" + ticket
	}

	return ticket
}
