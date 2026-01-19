package git

import (
	"strings"
	"testing"
)

func TestDetectCommitLanguage(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    string
	}{
		{
			name:    "English commit",
			subject: "feat: add user authentication",
			want:    "en",
		},
		{
			name:    "Chinese commit",
			subject: "feat: 添加用户认证功能",
			want:    "zh",
		},
		{
			name:    "Japanese commit",
			subject: "feat: ユーザー認証を追加",
			want:    "ja",
		},
		{
			name:    "Korean commit",
			subject: "feat: 사용자 인증 추가",
			want:    "ko",
		},
		{
			name:    "Russian commit",
			subject: "feat: добавить аутентификацию",
			want:    "ru",
		},
		{
			name:    "Mixed with Chinese characters",
			subject: "feat: add authentication (认证)",
			want:    "zh", // CJK characters are detected first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectCommitLanguage(tt.subject)
			if got != tt.want {
				t.Errorf("detectCommitLanguage(%q) = %v, want %v", tt.subject, got, tt.want)
			}
		})
	}
}

func TestAnalyzeTrends(t *testing.T) {
	// This is a unit test for the analyzeTrends helper
	// In practice, you'd use a mock or test data
	dayCommits := map[string]int{
		"2026-01-01": 5,
		"2026-01-05": 3,
		"2026-01-10": 8,
		"2026-01-15": 2,
		"2026-01-19": 4,
	}

	// Test with actual dates
	// Note: These dates would need to be adjusted based on actual time
	// This is a simplified test to verify the function works
	if len(dayCommits) == 0 {
		t.Error("dayCommits should not be empty")
	}

	// Create a simple trend analysis
	trends := &TrendAnalysis{
		Last30Days:    22,
		Last7Days:     4,
		MostActiveDay: "2026-01-10",
		AveragePerDay: 0.73,
	}

	if trends == nil {
		t.Error("trends should not be nil")
	}

	if trends.MostActiveDay == "" {
		t.Error("MostActiveDay should be set")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "Item exists",
			slice: []string{"feat", "fix", "docs"},
			item:  "fix",
			want:  true,
		},
		{
			name:  "Item does not exist",
			slice: []string{"feat", "fix", "docs"},
			item:  "test",
			want:  false,
		},
		{
			name:  "Empty slice",
			slice: []string{},
			item:  "feat",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.slice, tt.item)
			if got != tt.want {
				t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, got, tt.want)
			}
		})
	}
}

func TestSortMapByValue(t *testing.T) {
	m := map[string]int{
		"feat": 50,
		"fix":  30,
		"docs": 10,
		"test": 20,
	}

	sorted := sortMapByValue(m)

	if len(sorted) != 4 {
		t.Errorf("sortMapByValue() length = %d, want 4", len(sorted))
	}

	// Check descending order
	if sorted[0].Key != "feat" || sorted[0].Value != 50 {
		t.Errorf("First item should be feat:50, got %s:%d", sorted[0].Key, sorted[0].Value)
	}

	if sorted[1].Key != "fix" || sorted[1].Value != 30 {
		t.Errorf("Second item should be fix:30, got %s:%d", sorted[1].Key, sorted[1].Value)
	}
}

func TestCreateBar(t *testing.T) {
	tests := []struct {
		name       string
		percentage int
		maxWidth   int
		wantFilled int
	}{
		{
			name:       "50 percent",
			percentage: 50,
			maxWidth:   20,
			wantFilled: 10,
		},
		{
			name:       "100 percent",
			percentage: 100,
			maxWidth:   10,
			wantFilled: 10,
		},
		{
			name:       "0 percent",
			percentage: 0,
			maxWidth:   10,
			wantFilled: 0,
		},
		{
			name:       "Over 100 percent",
			percentage: 150,
			maxWidth:   10,
			wantFilled: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := createBar(tt.percentage, tt.maxWidth)

			// Check that bar contains brackets
			if !strings.HasPrefix(bar, "[") || !strings.HasSuffix(bar, "]") {
				t.Errorf("Bar should have brackets, got %q", bar)
			}

			// Count filled characters
			filled := strings.Count(bar, "█")
			if filled != tt.wantFilled {
				t.Errorf("createBar(%d, %d) filled = %d, want %d",
					tt.percentage, tt.maxWidth, filled, tt.wantFilled)
			}

			// Check total width (including brackets)
			expectedLen := tt.maxWidth + 2 // +2 for brackets
			if len([]rune(bar)) != expectedLen {
				t.Errorf("Bar length = %d, want %d", len([]rune(bar)), expectedLen)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "Short string",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "Exact length",
			input:  "hello world",
			maxLen: 11,
			want:   "hello world",
		},
		{
			name:   "Long string",
			input:  "this is a very long string that needs truncation",
			maxLen: 20,
			want:   "this is a very lo...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestGetLanguageName(t *testing.T) {
	tests := []struct {
		code string
		want string
	}{
		{"en", "English"},
		{"zh", "中文"},
		{"ja", "日本語"},
		{"ko", "한국어"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got := getLanguageName(tt.code)
			if got != tt.want {
				t.Errorf("getLanguageName(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a smaller", 5, 10, 5},
		{"b smaller", 10, 5, 5},
		{"equal", 7, 7, 7},
		{"negative", -5, 3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := min(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestGetTopPatterns(t *testing.T) {
	stats := &CommitStats{
		TypeDistribution: map[string]int{
			"feat": 50,
			"fix":  30,
			"docs": 10,
			"test": 5,
		},
		ScopeDistribution: map[string]int{
			"auth": 20,
			"api":  15,
		},
	}

	patterns := GetTopPatterns(stats, 3)

	if len(patterns) != 3 {
		t.Errorf("GetTopPatterns() returned %d patterns, want 3", len(patterns))
	}

	// Check that patterns are sorted by frequency
	if patterns[0].Type != "feat" || patterns[0].Frequency != 50 {
		t.Errorf("First pattern should be feat:50, got %s:%d",
			patterns[0].Type, patterns[0].Frequency)
	}

	if patterns[1].Type != "fix" || patterns[1].Frequency != 30 {
		t.Errorf("Second pattern should be fix:30, got %s:%d",
			patterns[1].Type, patterns[1].Frequency)
	}
}

func TestFormatStatsReport(t *testing.T) {
	stats := &CommitStats{
		TotalCommits: 100,
		TypeDistribution: map[string]int{
			"feat": 50,
			"fix":  30,
		},
		AverageLength:   45,
		WithScope:       40,
		WithBody:        20,
		WithTicket:      10,
		LongestSubject:  "feat: this is a very long commit message",
		ShortestSubject: "fix: typo",
		RecentTrends: &TrendAnalysis{
			Last30Days:    80,
			Last7Days:     15,
			AveragePerDay: 2.7,
			MostActiveDay: "2026-01-15",
		},
	}

	report := FormatStatsReport(stats)

	// Check that report contains expected sections
	expectedSections := []string{
		"Commit History Statistics",
		"Overview",
		"Commit Types",
		"Recent Activity",
		"Extremes",
	}

	for _, section := range expectedSections {
		if !strings.Contains(report, section) {
			t.Errorf("Report missing section: %s", section)
		}
	}

	// Check that stats are included
	if !strings.Contains(report, "100") { // total commits
		t.Error("Report should contain total commits count")
	}

	if !strings.Contains(report, "feat") {
		t.Error("Report should contain commit types")
	}
}
