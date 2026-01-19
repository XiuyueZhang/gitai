package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

// CommitStats holds statistics about commit history
type CommitStats struct {
	TotalCommits      int
	TypeDistribution  map[string]int // feat: 45, fix: 30, etc.
	ScopeDistribution map[string]int // auth: 10, api: 8, etc.
	AuthorStats       map[string]int // author -> commit count
	AverageLength     int            // average subject line length
	LongestSubject    string
	ShortestSubject   string
	WithScope         int // commits with scope
	WithBody          int // commits with body
	WithTicket        int // commits with ticket numbers
	CommonVerbs       map[string]int // add: 50, update: 30, etc.
	LanguageUsage     map[string]int // en: 80, zh: 20, etc.
	TimeDistribution  map[string]int // hour of day -> commit count
	DayDistribution   map[string]int // weekday -> commit count
	RecentTrends      *TrendAnalysis
}

// TrendAnalysis shows recent trends
type TrendAnalysis struct {
	Last30Days    int
	Last7Days     int
	LastWeek      int
	MostActiveDay string
	AveragePerDay float64
}

// CommitPattern represents a learned pattern from history
type CommitPattern struct {
	Type           string
	Scope          string
	CommonPhrases  []string
	ExampleMessage string
	Frequency      int
}

// AnalyzeCommitHistory analyzes commit history and returns statistics
func AnalyzeCommitHistory(limit int) (*CommitStats, error) {
	stats := &CommitStats{
		TypeDistribution:  make(map[string]int),
		ScopeDistribution: make(map[string]int),
		AuthorStats:       make(map[string]int),
		CommonVerbs:       make(map[string]int),
		LanguageUsage:     make(map[string]int),
		TimeDistribution:  make(map[string]int),
		DayDistribution:   make(map[string]int),
	}

	// Get commits with detailed format
	// Format: %H|%an|%s|%b|%ad
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", limit),
		"--pretty=format:%H|%an|%s|%b|%ad", "--date=iso")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(commits) == 1 && commits[0] == "" {
		return stats, nil
	}

	stats.TotalCommits = len(commits)

	// Patterns
	typePattern := regexp.MustCompile(`^(\w+)(\(([^)]+)\))?:\s*(.+)$`)
	ticketPattern := regexp.MustCompile(`\[([\w-]+)\]|\b([A-Z]+-\d+)\b`)
	verbPattern := regexp.MustCompile(`^(\w+)(\(([^)]+)\))?:\s*(\w+)`)

	var totalLength int
	longest := ""
	shortest := commits[0]

	now := time.Now()
	last30Days := now.AddDate(0, 0, -30)
	last7Days := now.AddDate(0, 0, -7)
	dayCommits := make(map[string]int)

	for _, commit := range commits {
		parts := strings.Split(commit, "|")
		if len(parts) < 5 {
			continue
		}

		// hash := parts[0]
		author := parts[1]
		subject := parts[2]
		body := parts[3]
		dateStr := parts[4]

		// Author stats
		stats.AuthorStats[author]++

		// Subject length
		subjectLen := len(subject)
		totalLength += subjectLen
		if subjectLen > len(longest) {
			longest = subject
		}
		if subjectLen < len(shortest) {
			shortest = subject
		}

		// Parse date
		commitDate, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr)
		if err == nil {
			// Time distribution (hour of day)
			hour := fmt.Sprintf("%02d:00", commitDate.Hour())
			stats.TimeDistribution[hour]++

			// Day distribution
			weekday := commitDate.Weekday().String()
			stats.DayDistribution[weekday]++

			// Trends
			dateKey := commitDate.Format("2006-01-02")
			dayCommits[dateKey]++
		}

		// Type and scope analysis
		if matches := typePattern.FindStringSubmatch(subject); matches != nil {
			commitType := matches[1]
			scope := matches[3]
			// message := matches[4]

			stats.TypeDistribution[commitType]++

			if scope != "" {
				stats.ScopeDistribution[scope]++
				stats.WithScope++
			}
		}

		// Body analysis
		if strings.TrimSpace(body) != "" {
			stats.WithBody++
		}

		// Ticket detection
		if ticketPattern.MatchString(subject) {
			stats.WithTicket++
		}

		// Verb extraction
		if matches := verbPattern.FindStringSubmatch(subject); matches != nil && len(matches) > 4 {
			verb := strings.ToLower(matches[4])
			stats.CommonVerbs[verb]++
		}

		// Language detection
		lang := detectCommitLanguage(subject)
		stats.LanguageUsage[lang]++
	}

	// Calculate averages
	if stats.TotalCommits > 0 {
		stats.AverageLength = totalLength / stats.TotalCommits
	}
	stats.LongestSubject = longest
	stats.ShortestSubject = shortest

	// Analyze trends
	stats.RecentTrends = analyzeTrends(dayCommits, now, last30Days, last7Days)

	return stats, nil
}

// analyzeTrends analyzes recent commit trends
func analyzeTrends(dayCommits map[string]int, now, last30Days, last7Days time.Time) *TrendAnalysis {
	trends := &TrendAnalysis{}

	var total30, total7 int
	maxDay := ""
	maxCount := 0

	for dateStr, count := range dayCommits {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if date.After(last30Days) {
			total30 += count
		}
		if date.After(last7Days) {
			total7 += count
		}

		if count > maxCount {
			maxCount = count
			maxDay = dateStr
		}
	}

	trends.Last30Days = total30
	trends.Last7Days = total7
	trends.MostActiveDay = maxDay

	if total30 > 0 {
		trends.AveragePerDay = float64(total30) / 30.0
	}

	return trends
}

// detectCommitLanguage detects the language of a commit message
func detectCommitLanguage(subject string) string {
	// Simple heuristic - check for CJK characters
	for _, r := range subject {
		if r >= 0x4E00 && r <= 0x9FFF {
			return "zh"
		}
		if (r >= 0x3040 && r <= 0x309F) || (r >= 0x30A0 && r <= 0x30FF) {
			return "ja"
		}
		if r >= 0xAC00 && r <= 0xD7AF {
			return "ko"
		}
		if r >= 0x0400 && r <= 0x04FF {
			return "ru"
		}
	}
	return "en"
}

// FindSimilarCommits finds commits similar to the current changes
func FindSimilarCommits(changedFiles []string, limit int) ([]string, error) {
	if len(changedFiles) == 0 {
		return []string{}, nil
	}

	// Build a pattern to search for commits affecting similar files
	var filePatterns []string
	for _, file := range changedFiles {
		// Get directory or base name for similarity
		parts := strings.Split(file, "/")
		if len(parts) > 1 {
			filePatterns = append(filePatterns, parts[len(parts)-2]) // directory name
		}
		filePatterns = append(filePatterns, parts[len(parts)-1]) // file name
	}

	// Get commits that modified similar files
	var similarCommits []string
	for _, pattern := range filePatterns {
		cmd := exec.Command("git", "log", fmt.Sprintf("-%d", limit),
			"--pretty=format:%s", "--", "*"+pattern+"*")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		commits := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, commit := range commits {
			if commit != "" && !contains(similarCommits, commit) {
				similarCommits = append(similarCommits, commit)
			}
		}

		if len(similarCommits) >= limit {
			break
		}
	}

	// Limit results
	if len(similarCommits) > limit {
		similarCommits = similarCommits[:limit]
	}

	return similarCommits, nil
}

// GetTopPatterns extracts the most common commit patterns
func GetTopPatterns(stats *CommitStats, topN int) []CommitPattern {
	var patterns []CommitPattern

	// Sort types by frequency
	type kv struct {
		Key   string
		Value int
	}

	var typeList []kv
	for k, v := range stats.TypeDistribution {
		typeList = append(typeList, kv{k, v})
	}
	sort.Slice(typeList, func(i, j int) bool {
		return typeList[i].Value > typeList[j].Value
	})

	// Create patterns for top types
	for i := 0; i < topN && i < len(typeList); i++ {
		commitType := typeList[i].Key
		pattern := CommitPattern{
			Type:      commitType,
			Frequency: typeList[i].Value,
		}

		// Find most common scope for this type
		// (This is simplified - in practice you'd correlate type+scope)
		if len(stats.ScopeDistribution) > 0 {
			var scopeList []kv
			for k, v := range stats.ScopeDistribution {
				scopeList = append(scopeList, kv{k, v})
			}
			sort.Slice(scopeList, func(i, j int) bool {
				return scopeList[i].Value > scopeList[j].Value
			})
			if len(scopeList) > 0 {
				pattern.Scope = scopeList[0].Key
			}
		}

		patterns = append(patterns, pattern)
	}

	return patterns
}

// contains checks if a string slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// FormatStatsReport generates a human-readable stats report
func FormatStatsReport(stats *CommitStats) string {
	var report strings.Builder

	report.WriteString("📊 Commit History Statistics\n")
	report.WriteString(strings.Repeat("=", 50) + "\n\n")

	// Overview
	report.WriteString("📈 Overview\n")
	report.WriteString(fmt.Sprintf("  Total Commits: %d\n", stats.TotalCommits))
	report.WriteString(fmt.Sprintf("  Average Subject Length: %d characters\n", stats.AverageLength))
	report.WriteString(fmt.Sprintf("  With Scope: %d (%.1f%%)\n", stats.WithScope,
		float64(stats.WithScope)/float64(stats.TotalCommits)*100))
	report.WriteString(fmt.Sprintf("  With Body: %d (%.1f%%)\n", stats.WithBody,
		float64(stats.WithBody)/float64(stats.TotalCommits)*100))
	report.WriteString(fmt.Sprintf("  With Ticket: %d (%.1f%%)\n\n", stats.WithTicket,
		float64(stats.WithTicket)/float64(stats.TotalCommits)*100))

	// Type distribution
	report.WriteString("🏷️  Commit Types\n")
	typeList := sortMapByValue(stats.TypeDistribution)
	for _, kv := range typeList[:min(10, len(typeList))] {
		percentage := float64(kv.Value) / float64(stats.TotalCommits) * 100
		bar := createBar(int(percentage), 20)
		report.WriteString(fmt.Sprintf("  %-12s %s %3d (%.1f%%)\n", kv.Key, bar, kv.Value, percentage))
	}
	report.WriteString("\n")

	// Scope distribution (if any)
	if len(stats.ScopeDistribution) > 0 {
		report.WriteString("📦 Top Scopes\n")
		scopeList := sortMapByValue(stats.ScopeDistribution)
		for _, kv := range scopeList[:min(8, len(scopeList))] {
			percentage := float64(kv.Value) / float64(stats.WithScope) * 100
			bar := createBar(int(percentage), 15)
			report.WriteString(fmt.Sprintf("  %-15s %s %3d\n", kv.Key, bar, kv.Value))
		}
		report.WriteString("\n")
	}

	// Common verbs
	if len(stats.CommonVerbs) > 0 {
		report.WriteString("🔤 Common Action Verbs\n")
		verbList := sortMapByValue(stats.CommonVerbs)
		for _, kv := range verbList[:min(8, len(verbList))] {
			report.WriteString(fmt.Sprintf("  %-12s %3d\n", kv.Key, kv.Value))
		}
		report.WriteString("\n")
	}

	// Language usage
	if len(stats.LanguageUsage) > 1 {
		report.WriteString("🌍 Language Usage\n")
		langList := sortMapByValue(stats.LanguageUsage)
		for _, kv := range langList {
			percentage := float64(kv.Value) / float64(stats.TotalCommits) * 100
			bar := createBar(int(percentage), 20)
			langName := getLanguageName(kv.Key)
			report.WriteString(fmt.Sprintf("  %-10s %s %.1f%%\n", langName, bar, percentage))
		}
		report.WriteString("\n")
	}

	// Time distribution
	if len(stats.TimeDistribution) > 0 {
		report.WriteString("⏰ Commit Time Distribution\n")
		timeList := sortMapByValue(stats.TimeDistribution)
		for _, kv := range timeList[:min(5, len(timeList))] {
			bar := createBar(kv.Value*2, 15)
			report.WriteString(fmt.Sprintf("  %s  %s %d\n", kv.Key, bar, kv.Value))
		}
		report.WriteString("\n")
	}

	// Day distribution
	if len(stats.DayDistribution) > 0 {
		report.WriteString("📅 Commit Day Distribution\n")
		// Sort by weekday order
		weekdays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		maxCount := 0
		for _, count := range stats.DayDistribution {
			if count > maxCount {
				maxCount = count
			}
		}
		for _, day := range weekdays {
			if count, ok := stats.DayDistribution[day]; ok {
				percentage := float64(count) / float64(maxCount) * 100
				bar := createBar(int(percentage/5), 20)
				report.WriteString(fmt.Sprintf("  %-10s %s %3d\n", day, bar, count))
			}
		}
		report.WriteString("\n")
	}

	// Recent trends
	if stats.RecentTrends != nil {
		report.WriteString("📊 Recent Activity\n")
		report.WriteString(fmt.Sprintf("  Last 30 days: %d commits (%.1f/day avg)\n",
			stats.RecentTrends.Last30Days, stats.RecentTrends.AveragePerDay))
		report.WriteString(fmt.Sprintf("  Last 7 days:  %d commits\n", stats.RecentTrends.Last7Days))
		if stats.RecentTrends.MostActiveDay != "" {
			report.WriteString(fmt.Sprintf("  Most active:  %s\n", stats.RecentTrends.MostActiveDay))
		}
		report.WriteString("\n")
	}

	// Authors (top contributors)
	if len(stats.AuthorStats) > 0 {
		report.WriteString("👥 Top Contributors\n")
		authorList := sortMapByValue(stats.AuthorStats)
		for _, kv := range authorList[:min(5, len(authorList))] {
			percentage := float64(kv.Value) / float64(stats.TotalCommits) * 100
			bar := createBar(int(percentage), 15)
			report.WriteString(fmt.Sprintf("  %-25s %s %.1f%%\n", kv.Key, bar, percentage))
		}
		report.WriteString("\n")
	}

	// Extremes
	report.WriteString("📏 Extremes\n")
	report.WriteString(fmt.Sprintf("  Longest:  %s\n", truncate(stats.LongestSubject, 60)))
	report.WriteString(fmt.Sprintf("  Shortest: %s\n", truncate(stats.ShortestSubject, 60)))

	return report.String()
}

// Helper functions

type kv struct {
	Key   string
	Value int
}

func sortMapByValue(m map[string]int) []kv {
	var list []kv
	for k, v := range m {
		list = append(list, kv{k, v})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Value > list[j].Value
	})
	return list
}

func createBar(percentage, maxWidth int) string {
	if percentage > 100 {
		percentage = 100
	}
	width := percentage * maxWidth / 100
	if width < 0 {
		width = 0
	}
	return "[" + strings.Repeat("█", width) + strings.Repeat("░", maxWidth-width) + "]"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getLanguageName(code string) string {
	names := map[string]string{
		"en": "English",
		"zh": "中文",
		"ja": "日本語",
		"ko": "한국어",
		"de": "Deutsch",
		"fr": "Français",
		"es": "Español",
		"pt": "Português",
		"ru": "Русский",
		"it": "Italiano",
	}
	if name, ok := names[code]; ok {
		return name
	}
	return code
}
