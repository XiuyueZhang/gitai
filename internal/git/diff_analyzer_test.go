package git

import (
	"strings"
	"testing"
)

func TestAnalyzeDiff(t *testing.T) {
	tests := []struct {
		name        string
		diff        string
		maxLength   int
		wantFiles   int
		wantComplex string
	}{
		{
			name: "Simple Go file change",
			diff: `diff --git a/main.go b/main.go
index 123..456 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,7 @@
 package main

+import "fmt"
+
 func main() {
-	println("hello")
+	fmt.Println("hello world")
 }`,
			maxLength:   2000,
			wantFiles:   1,
			wantComplex: "simple",
		},
		{
			name: "Function addition",
			diff: `diff --git a/pkg/auth.go b/pkg/auth.go
index 123..456 100644
--- a/pkg/auth.go
+++ b/pkg/auth.go
@@ -1,5 +1,15 @@
 package pkg

+func NewAuthenticator() *Authenticator {
+	return &Authenticator{}
+}
+
+func (a *Authenticator) Login(user string) error {
+	return nil
+}`,
			maxLength:   3000,
			wantFiles:   1,
			wantComplex: "simple",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeDiff(tt.diff, tt.maxLength)

			if len(analysis.FileSummaries) != tt.wantFiles {
				t.Errorf("AnalyzeDiff() files = %d, want %d", len(analysis.FileSummaries), tt.wantFiles)
			}

			if analysis.ChangeComplexity != tt.wantComplex {
				t.Errorf("AnalyzeDiff() complexity = %s, want %s", analysis.ChangeComplexity, tt.wantComplex)
			}

			if analysis.SmartDiff == "" {
				t.Errorf("AnalyzeDiff() SmartDiff is empty")
			}

			// Log extracted information for debugging
			t.Logf("Extracted %d key changes, %d import changes",
				len(analysis.KeyChanges), len(analysis.ImportChanges))
		})
	}
}

func TestAnalyzeFileType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"main.go", "go"},
		{"app.js", "javascript"},
		{"component.tsx", "typescript"},
		{"script.py", "python"},
		{"Main.java", "java"},
		{"config.yaml", "yaml"},
		{"package.json", "json"},
		{"README.md", "markdown"},
		{"unknown.xyz", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := analyzeFileType(tt.path); got != tt.want {
				t.Errorf("analyzeFileType(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"main_test.go", true},
		{"app.test.js", true},
		{"src/test/auth.go", true},
		{"component.spec.ts", true},
		{"__tests__/app.js", true},
		{"main.go", false},
		{"src/auth.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTestFile(tt.path); got != tt.want {
				t.Errorf("isTestFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsConfigFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"package.json", true},
		{"go.mod", true},
		{"config.yaml", true},
		{".env", true},
		{"Dockerfile", true},
		{"Makefile", true},
		{"main.go", false},
		{"src/app.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isConfigFile(tt.path); got != tt.want {
				t.Errorf("isConfigFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestExtractGoChanges(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int // number of changes expected
	}{
		{
			name: "Function definition",
			line: "func NewAuthenticator() *Authenticator {",
			want: 1,
		},
		{
			name: "Method definition",
			line: "func (a *Auth) Login(user string) error {",
			want: 1,
		},
		{
			name: "Type definition",
			line: "type User struct {",
			want: 1,
		},
		{
			name: "Interface definition",
			line: "type Handler interface {",
			want: 1,
		},
		{
			name: "Regular code",
			line: "return nil",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractGoChanges(tt.line)
			if len(got) != tt.want {
				t.Errorf("extractGoChanges(%q) = %d changes, want %d", tt.line, len(got), tt.want)
			}
		})
	}
}

func TestExtractImportChanges(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,3 +1,5 @@
 package main

+import "fmt"
+import "net/http"
+
 func main() {`

	changes := extractImportChanges(diff)

	if len(changes) < 2 {
		t.Errorf("extractImportChanges() = %d imports, want at least 2", len(changes))
	}

	// Check that import statements are captured
	hasImport := false
	for _, change := range changes {
		if strings.Contains(change, "import") {
			hasImport = true
			break
		}
	}

	if !hasImport {
		t.Errorf("extractImportChanges() did not capture import statements")
	}
}

func TestDetermineComplexity(t *testing.T) {
	tests := []struct {
		name           string
		totalAdditions int
		totalDeletions int
		modifiedFiles  int
		want           string
	}{
		{"Simple change", 10, 5, 1, "simple"},
		{"Moderate change", 150, 50, 5, "moderate"},
		{"Complex change", 600, 400, 12, "complex"},
		{"Large file count", 50, 50, 15, "complex"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &DiffAnalysis{
				TotalAdditions: tt.totalAdditions,
				TotalDeletions: tt.totalDeletions,
				ModifiedFiles:  tt.modifiedFiles,
			}

			if got := determineComplexity(analysis); got != tt.want {
				t.Errorf("determineComplexity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitDiffByFile(t *testing.T) {
	diff := `diff --git a/file1.go b/file1.go
index 123..456
--- a/file1.go
+++ b/file1.go
@@ -1 +1 @@
-old content
+new content
diff --git a/file2.go b/file2.go
index 789..abc
--- a/file2.go
+++ b/file2.go
@@ -1 +1 @@
-old content 2
+new content 2`

	files := splitDiffByFile(diff)

	if len(files) != 2 {
		t.Errorf("splitDiffByFile() = %d files, want 2", len(files))
	}

	// Check that each file diff contains the expected marker
	for i, fileDiff := range files {
		if !strings.Contains(fileDiff, "diff --git") {
			t.Errorf("File diff %d does not contain 'diff --git' marker", i)
		}
	}
}

func TestGenerateSmartDiff(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
index 123..456
--- a/main.go
+++ b/main.go
@@ -1,5 +1,10 @@
 package main

+import "fmt"
+
 func main() {
-	println("hello")
+	fmt.Println("hello world")
+}
+
+func NewFunction() {
+	// New functionality
 }`

	files := splitDiffByFile(diff)
	analysis := AnalyzeDiff(diff, 2000)

	smartDiff := generateSmartDiff(diff, files, analysis, 500)

	if smartDiff == "" {
		t.Error("generateSmartDiff() returned empty string")
	}

	// Check that summary is included
	if !strings.Contains(smartDiff, "DIFF SUMMARY") {
		t.Error("generateSmartDiff() does not contain summary header")
	}

	// Check that file information is included
	if !strings.Contains(smartDiff, "main.go") {
		t.Error("generateSmartDiff() does not contain file name")
	}
}

func TestUniqueStrings(t *testing.T) {
	input := []string{"a", "b", "a", "c", "b", ""}
	expected := []string{"a", "b", "c"}

	result := uniqueStrings(input)

	if len(result) != len(expected) {
		t.Errorf("uniqueStrings() length = %d, want %d", len(result), len(expected))
	}

	// Check that all unique values are present
	for _, want := range expected {
		found := false
		for _, got := range result {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("uniqueStrings() missing value %q", want)
		}
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is a ..."},
		{"exactly12ch", 11, "exactly12ch"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}
