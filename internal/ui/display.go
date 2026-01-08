package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/xyue92/gitai/internal/git"
)

// Display provides formatted output functions
type Display struct {
	NoColor bool
}

// NewDisplay creates a new Display instance
func NewDisplay() *Display {
	return &Display{}
}

// ShowHeader displays the application header
func (d *Display) ShowHeader() {
	title := color.New(color.FgCyan, color.Bold)
	if d.NoColor {
		fmt.Println("Git Commit AI Assistant")
	} else {
		title.Println("📝 Git Commit AI Assistant")
	}
	fmt.Println()
}

// ShowChangedFiles displays the list of changed files
func (d *Display) ShowChangedFiles(files []git.FileChange) {
	if len(files) == 0 {
		return
	}

	bold := color.New(color.Bold)
	bold.Printf("Changed files (%d):\n", len(files))

	for _, file := range files {
		additions := file.Additions
		deletions := file.Deletions

		if additions == "-" {
			additions = "bin"
		}
		if deletions == "-" {
			deletions = "bin"
		}

		green := color.New(color.FgGreen)
		red := color.New(color.FgRed)

		if d.NoColor {
			fmt.Printf("  ✓ %s (+%s, -%s)\n", file.File, additions, deletions)
		} else {
			fmt.Printf("  ✓ %s (", file.File)
			green.Printf("+%s", additions)
			fmt.Printf(", ")
			red.Printf("-%s", deletions)
			fmt.Println(")")
		}
	}
	fmt.Println()
}

// ShowGenerating displays a "generating..." message
func (d *Display) ShowGenerating() {
	yellow := color.New(color.FgYellow)
	if d.NoColor {
		fmt.Println("Generating commit message...")
	} else {
		yellow.Println("🤖 Generating commit message...")
	}
	fmt.Println()
}

// ShowCommitMessage displays the generated commit message in a box
func (d *Display) ShowCommitMessage(message string) {
	bold := color.New(color.Bold)
	bold.Println("Generated message:")

	// Create a simple box around the message
	lines := strings.Split(message, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Ensure minimum width
	if maxLen < 40 {
		maxLen = 40
	}

	// Top border
	fmt.Println("┌─" + strings.Repeat("─", maxLen) + "─┐")

	// Content
	for _, line := range lines {
		padding := maxLen - len(line)
		cyan := color.New(color.FgCyan)
		if d.NoColor {
			fmt.Printf("│ %s%s │\n", line, strings.Repeat(" ", padding))
		} else {
			fmt.Print("│ ")
			cyan.Print(line)
			fmt.Printf("%s │\n", strings.Repeat(" ", padding))
		}
	}

	// Bottom border
	fmt.Println("└─" + strings.Repeat("─", maxLen) + "─┘")
	fmt.Println()
}

// ShowSuccess displays a success message
func (d *Display) ShowSuccess(message string) {
	green := color.New(color.FgGreen, color.Bold)
	if d.NoColor {
		fmt.Println("✓ " + message)
	} else {
		green.Println("✨ " + message)
	}
}

// ShowError displays an error message
func (d *Display) ShowError(err error) {
	red := color.New(color.FgRed, color.Bold)
	if d.NoColor {
		fmt.Printf("✗ Error: %v\n", err)
	} else {
		red.Print("❌ Error: ")
		fmt.Println(err)
	}
}

// ShowWarning displays a warning message
func (d *Display) ShowWarning(message string) {
	yellow := color.New(color.FgYellow, color.Bold)
	if d.NoColor {
		fmt.Println("⚠ " + message)
	} else {
		yellow.Println("⚠️  " + message)
	}
}

// ShowInfo displays an info message
func (d *Display) ShowInfo(message string) {
	blue := color.New(color.FgBlue)
	if d.NoColor {
		fmt.Println(message)
	} else {
		blue.Println(message)
	}
}

// ShowDryRun displays dry-run mode header
func (d *Display) ShowDryRun() {
	cyan := color.New(color.FgCyan, color.Bold)
	if d.NoColor {
		fmt.Println("Dry-run mode - no commit will be created")
	} else {
		cyan.Println("🔍 Dry-run mode - no commit will be created")
	}
	fmt.Println()
}

// ShowCommitSuccess displays success after commit
func (d *Display) ShowCommitSuccess(message string, files []string) {
	d.ShowSuccess("Commit created successfully!")
	fmt.Println()

	bold := color.New(color.Bold)
	bold.Println("Commit message:")
	fmt.Println(message)
	fmt.Println()

	if len(files) > 0 {
		bold.Println("Files changed:")
		for _, file := range files {
			fmt.Printf("  %s\n", file)
		}
		fmt.Println()
	}

	d.ShowInfo("View commit: git show HEAD")
}
