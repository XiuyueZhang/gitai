package git

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// DiffAnalysis contains enhanced analysis of git diff
type DiffAnalysis struct {
	FileSummaries    []FileSummary // Per-file change summaries
	SmartDiff        string        // Intelligently truncated diff
	TotalAdditions   int           // Total lines added
	TotalDeletions   int           // Total lines deleted
	ModifiedFiles    int           // Number of files modified
	KeyChanges       []string      // Key code changes (functions, classes, etc.)
	ImportChanges    []string      // Import/dependency changes
	IsLargeChange    bool          // Whether this is a large refactoring
	ChangeComplexity string        // simple, moderate, complex
}

// FileSummary represents changes in a single file
type FileSummary struct {
	Path           string
	Status         string // modified, added, deleted, renamed
	Additions      int
	Deletions      int
	FileType       string // go, js, ts, py, etc.
	IsTestFile     bool
	IsConfigFile   bool
	KeyChanges     []string // Important changes in this file
}

// AnalyzeDiff performs intelligent analysis of git diff
func AnalyzeDiff(fullDiff string, maxLength int) *DiffAnalysis {
	analysis := &DiffAnalysis{
		FileSummaries: make([]FileSummary, 0),
		KeyChanges:    make([]string, 0),
		ImportChanges: make([]string, 0),
	}

	if fullDiff == "" {
		return analysis
	}

	// Split diff by files
	files := splitDiffByFile(fullDiff)

	// Analyze each file
	for _, fileDiff := range files {
		summary := analyzeFileDiff(fileDiff)
		analysis.FileSummaries = append(analysis.FileSummaries, summary)

		analysis.TotalAdditions += summary.Additions
		analysis.TotalDeletions += summary.Deletions
		analysis.ModifiedFiles++

		// Collect key changes
		analysis.KeyChanges = append(analysis.KeyChanges, summary.KeyChanges...)

		// Collect import changes
		imports := extractImportChanges(fileDiff)
		analysis.ImportChanges = append(analysis.ImportChanges, imports...)
	}

	// Determine change complexity
	analysis.ChangeComplexity = determineComplexity(analysis)
	analysis.IsLargeChange = analysis.TotalAdditions+analysis.TotalDeletions > 500

	// Generate smart diff (intelligently truncated)
	analysis.SmartDiff = generateSmartDiff(fullDiff, files, analysis, maxLength)

	return analysis
}

// splitDiffByFile splits the full diff into per-file diffs
func splitDiffByFile(diff string) []string {
	// Split by "diff --git" marker
	parts := regexp.MustCompile(`(?m)^diff --git`).Split(diff, -1)

	files := make([]string, 0)
	for i, part := range parts {
		if i == 0 && part == "" {
			continue // Skip empty first part
		}
		// Re-add the "diff --git" marker
		if i > 0 {
			part = "diff --git" + part
		}
		files = append(files, part)
	}

	return files
}

// analyzeFileDiff analyzes a single file's diff
func analyzeFileDiff(fileDiff string) FileSummary {
	summary := FileSummary{
		KeyChanges: make([]string, 0),
	}

	lines := strings.Split(fileDiff, "\n")

	// Extract file path and status
	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Extract file path: diff --git a/path b/path
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				summary.Path = strings.TrimPrefix(parts[3], "b/")
			}
		} else if strings.HasPrefix(line, "new file") {
			summary.Status = "added"
		} else if strings.HasPrefix(line, "deleted file") {
			summary.Status = "deleted"
		} else if strings.HasPrefix(line, "rename from") {
			summary.Status = "renamed"
		} else if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			if summary.Status == "" {
				summary.Status = "modified"
			}
		}

		// Count additions and deletions
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			summary.Additions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			summary.Deletions++
		}
	}

	// Determine file type
	summary.FileType = analyzeFileType(summary.Path)
	summary.IsTestFile = isTestFile(summary.Path)
	summary.IsConfigFile = isConfigFile(summary.Path)

	// Extract key changes (function definitions, class definitions, etc.)
	summary.KeyChanges = extractKeyChanges(fileDiff, summary.FileType)

	return summary
}

// analyzeFileType determines the file type from path
func analyzeFileType(path string) string {
	ext := ""
	if idx := strings.LastIndex(path, "."); idx != -1 {
		ext = path[idx+1:]
	}

	typeMap := map[string]string{
		"go":   "go",
		"js":   "javascript",
		"jsx":  "javascript",
		"ts":   "typescript",
		"tsx":  "typescript",
		"py":   "python",
		"java": "java",
		"rb":   "ruby",
		"rs":   "rust",
		"c":    "c",
		"cpp":  "cpp",
		"h":    "c",
		"hpp":  "cpp",
		"cs":   "csharp",
		"php":  "php",
		"swift": "swift",
		"kt":   "kotlin",
		"yaml": "yaml",
		"yml":  "yaml",
		"json": "json",
		"toml": "toml",
		"md":   "markdown",
		"sql":  "sql",
	}

	if t, ok := typeMap[ext]; ok {
		return t
	}
	return "unknown"
}

// isTestFile checks if the file is a test file
func isTestFile(path string) bool {
	lowerPath := strings.ToLower(path)
	return strings.Contains(lowerPath, "_test.") ||
		strings.Contains(lowerPath, ".test.") ||
		strings.Contains(lowerPath, "/test/") ||
		strings.Contains(lowerPath, "/tests/") ||
		strings.Contains(lowerPath, "__tests__") ||
		strings.Contains(lowerPath, ".spec.")
}

// isConfigFile checks if the file is a configuration file
func isConfigFile(path string) bool {
	lowerPath := strings.ToLower(path)
	configFiles := []string{
		"package.json", "go.mod", "go.sum", "cargo.toml",
		"requirements.txt", "gemfile", "pom.xml", "build.gradle",
		".yml", ".yaml", ".toml", ".json", ".config", ".env",
		"dockerfile", "makefile", ".gitignore", ".dockerignore",
	}

	for _, cf := range configFiles {
		if strings.Contains(lowerPath, cf) {
			return true
		}
	}
	return false
}

// extractKeyChanges extracts important code changes
func extractKeyChanges(fileDiff string, fileType string) []string {
	changes := make([]string, 0)
	lines := strings.Split(fileDiff, "\n")

	for _, line := range lines {
		// Only look at added/modified lines
		if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
			continue
		}

		trimmed := strings.TrimPrefix(line, "+")
		trimmed = strings.TrimSpace(trimmed)

		// Extract based on file type
		switch fileType {
		case "go":
			changes = append(changes, extractGoChanges(trimmed)...)
		case "javascript", "typescript":
			changes = append(changes, extractJSChanges(trimmed)...)
		case "python":
			changes = append(changes, extractPythonChanges(trimmed)...)
		default:
			changes = append(changes, extractGenericChanges(trimmed)...)
		}
	}

	// Remove duplicates and limit to top 5
	changes = uniqueStrings(changes)
	if len(changes) > 5 {
		changes = changes[:5]
	}

	return changes
}

// extractGoChanges extracts Go-specific changes
func extractGoChanges(line string) []string {
	changes := make([]string, 0)

	// Function definitions
	if matched, _ := regexp.MatchString(`^func\s+(\w+|\(\w+\s+\*?\w+\))\s+\w+`, line); matched {
		changes = append(changes, extractFunctionName(line, "go"))
	}

	// Type definitions
	if matched, _ := regexp.MatchString(`^type\s+\w+\s+(struct|interface)`, line); matched {
		changes = append(changes, extractTypeName(line))
	}

	return changes
}

// extractJSChanges extracts JavaScript/TypeScript changes
func extractJSChanges(line string) []string {
	changes := make([]string, 0)

	// Function definitions
	patterns := []string{
		`^function\s+\w+`,
		`^const\s+\w+\s*=\s*(async\s+)?\(`,
		`^export\s+(async\s+)?function\s+\w+`,
		`^\w+\s*\([^)]*\)\s*{`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			changes = append(changes, extractFunctionName(line, "js"))
			break
		}
	}

	// Class definitions
	if matched, _ := regexp.MatchString(`^(export\s+)?(default\s+)?class\s+\w+`, line); matched {
		changes = append(changes, extractClassName(line))
	}

	return changes
}

// extractPythonChanges extracts Python-specific changes
func extractPythonChanges(line string) []string {
	changes := make([]string, 0)

	// Function/method definitions
	if matched, _ := regexp.MatchString(`^(async\s+)?def\s+\w+`, line); matched {
		changes = append(changes, extractFunctionName(line, "python"))
	}

	// Class definitions
	if matched, _ := regexp.MatchString(`^class\s+\w+`, line); matched {
		changes = append(changes, extractClassName(line))
	}

	return changes
}

// extractGenericChanges extracts generic code patterns
func extractGenericChanges(line string) []string {
	changes := make([]string, 0)

	// Generic function-like patterns
	if matched, _ := regexp.MatchString(`\w+\s*\([^)]*\)\s*{?`, line); matched {
		if len(line) < 100 { // Avoid overly long lines
			changes = append(changes, fmt.Sprintf("function: %s", truncateString(line, 50)))
		}
	}

	return changes
}

// extractFunctionName extracts function name from line
func extractFunctionName(line, lang string) string {
	re := regexp.MustCompile(`\b(func|function|def|fn)\b\s+(\w+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 2 {
		return fmt.Sprintf("function %s", matches[2])
	}
	return ""
}

// extractClassName extracts class name from line
func extractClassName(line string) string {
	re := regexp.MustCompile(`\bclass\b\s+(\w+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return fmt.Sprintf("class %s", matches[1])
	}
	return ""
}

// extractTypeName extracts type name from line
func extractTypeName(line string) string {
	re := regexp.MustCompile(`\btype\b\s+(\w+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return fmt.Sprintf("type %s", matches[1])
	}
	return ""
}

// extractImportChanges extracts import/dependency changes
func extractImportChanges(fileDiff string) []string {
	changes := make([]string, 0)
	lines := strings.Split(fileDiff, "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "+++") {
			continue
		}

		trimmed := strings.TrimPrefix(line, "+")
		trimmed = strings.TrimSpace(trimmed)

		// Match import patterns
		patterns := []string{
			`^import\s+`,           // Go, Python, JS
			`^from\s+\S+\s+import`, // Python
			`^require\(`,           // Node.js
			`^use\s+`,              // Rust, PHP
		}

		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString(pattern, trimmed); matched {
				changes = append(changes, truncateString(trimmed, 60))
				break
			}
		}
	}

	return uniqueStrings(changes)
}

// determineComplexity determines the complexity of changes
func determineComplexity(analysis *DiffAnalysis) string {
	totalChanges := analysis.TotalAdditions + analysis.TotalDeletions
	fileCount := analysis.ModifiedFiles

	if totalChanges > 500 || fileCount > 10 {
		return "complex"
	} else if totalChanges > 100 || fileCount > 3 {
		return "moderate"
	}
	return "simple"
}

// generateSmartDiff creates an intelligently truncated diff
func generateSmartDiff(fullDiff string, files []string, analysis *DiffAnalysis, maxLength int) string {
	if len(fullDiff) <= maxLength {
		return fullDiff
	}

	var sb strings.Builder

	// 1. Add summary header
	sb.WriteString(fmt.Sprintf("DIFF SUMMARY (%d files, +%d/-%d lines)\n",
		analysis.ModifiedFiles, analysis.TotalAdditions, analysis.TotalDeletions))
	sb.WriteString(strings.Repeat("=", 60) + "\n\n")

	// 2. Add per-file summaries
	for _, summary := range analysis.FileSummaries {
		sb.WriteString(fmt.Sprintf("📄 %s [%s] +%d/-%d\n",
			summary.Path, summary.Status, summary.Additions, summary.Deletions))
		if len(summary.KeyChanges) > 0 {
			sb.WriteString(fmt.Sprintf("   Key changes: %s\n", strings.Join(summary.KeyChanges, ", ")))
		}
	}
	sb.WriteString("\n")

	// 3. Add import changes if any
	if len(analysis.ImportChanges) > 0 {
		sb.WriteString("📦 Import changes:\n")
		for _, imp := range analysis.ImportChanges {
			sb.WriteString(fmt.Sprintf("   %s\n", imp))
		}
		sb.WriteString("\n")
	}

	// 4. Add selected diff chunks (prioritize important files)
	remainingLength := maxLength - sb.Len()
	sb.WriteString("SELECTED DIFF CHUNKS:\n")
	sb.WriteString(strings.Repeat("-", 60) + "\n\n")

	// Sort files by importance (non-test, non-config first)
	sortedSummaries := make([]FileSummary, len(analysis.FileSummaries))
	copy(sortedSummaries, analysis.FileSummaries)
	sort.Slice(sortedSummaries, func(i, j int) bool {
		// Prioritize: non-test > non-config > larger changes
		if sortedSummaries[i].IsTestFile != sortedSummaries[j].IsTestFile {
			return !sortedSummaries[i].IsTestFile
		}
		if sortedSummaries[i].IsConfigFile != sortedSummaries[j].IsConfigFile {
			return !sortedSummaries[i].IsConfigFile
		}
		return (sortedSummaries[i].Additions + sortedSummaries[i].Deletions) >
			(sortedSummaries[j].Additions + sortedSummaries[j].Deletions)
	})

	// Add chunks from most important files
	for i, summary := range sortedSummaries {
		if i >= len(files) {
			break
		}

		// Find the corresponding file diff
		var fileDiff string
		for _, fd := range files {
			if strings.Contains(fd, summary.Path) {
				fileDiff = fd
				break
			}
		}

		if fileDiff == "" {
			continue
		}

		// Extract important chunks from this file
		chunk := extractImportantChunks(fileDiff, remainingLength/max(len(sortedSummaries)-i, 1))
		if chunk != "" {
			sb.WriteString(chunk)
			sb.WriteString("\n")
			remainingLength -= len(chunk)
		}

		if remainingLength < 500 {
			break
		}
	}

	// 5. Add note if truncated
	if len(fullDiff) > maxLength {
		sb.WriteString(fmt.Sprintf("\n... (diff truncated: %d/%d chars shown)\n",
			sb.Len(), len(fullDiff)))
	}

	return sb.String()
}

// extractImportantChunks extracts the most important parts of a file diff
func extractImportantChunks(fileDiff string, maxChunkLength int) string {
	lines := strings.Split(fileDiff, "\n")
	var sb strings.Builder

	// Keep file header
	for i, line := range lines {
		if i < 10 && (strings.HasPrefix(line, "diff ") ||
			strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++")) {
			sb.WriteString(line + "\n")
		}
	}

	// Extract hunks with meaningful changes
	inHunk := false
	hunkLines := make([]string, 0)
	contextLines := 0

	for _, line := range lines {
		// Start of a new hunk
		if strings.HasPrefix(line, "@@") {
			// Save previous hunk if it had changes
			if len(hunkLines) > 0 {
				sb.WriteString(strings.Join(hunkLines, "\n") + "\n")
			}
			hunkLines = []string{line}
			inHunk = true
			contextLines = 0
			continue
		}

		if !inHunk {
			continue
		}

		// Count context lines (not + or -)
		if !strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "-") {
			contextLines++
			// Skip excessive context
			if contextLines > 3 {
				continue
			}
		} else {
			contextLines = 0
		}

		hunkLines = append(hunkLines, line)

		// Stop if we've collected enough
		if sb.Len()+len(strings.Join(hunkLines, "\n")) > maxChunkLength {
			break
		}
	}

	// Add last hunk
	if len(hunkLines) > 0 && sb.Len() < maxChunkLength {
		sb.WriteString(strings.Join(hunkLines, "\n") + "\n")
	}

	return sb.String()
}

// Helper functions

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)

	for _, item := range slice {
		if item != "" && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
