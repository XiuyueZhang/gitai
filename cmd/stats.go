package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xyue92/gitai/internal/git"
)

var (
	statsLimit  int
	statsExport string
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show commit history statistics and patterns",
	Long: `Analyze your commit history and display comprehensive statistics including:
- Commit type distribution (feat, fix, docs, etc.)
- Scope usage patterns
- Common action verbs
- Language distribution
- Time and day patterns
- Recent activity trends
- Top contributors

This helps you understand your team's commit patterns and improve consistency.`,
	Example: `  # Show stats for last 100 commits
  gitai stats

  # Analyze last 500 commits
  gitai stats --limit 500

  # Export stats to JSON
  gitai stats --export stats.json`,
	RunE: runStats,
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().IntVarP(&statsLimit, "limit", "n", 100,
		"Number of commits to analyze")
	statsCmd.Flags().StringVarP(&statsExport, "export", "e", "",
		"Export statistics to JSON file")
}

func runStats(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !git.IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	fmt.Printf("🔍 Analyzing last %d commits...\n\n", statsLimit)

	// Analyze commit history
	stats, err := git.AnalyzeCommitHistory(statsLimit)
	if err != nil {
		return fmt.Errorf("failed to analyze commits: %w", err)
	}

	if stats.TotalCommits == 0 {
		fmt.Println("No commits found in the repository.")
		return nil
	}

	// Display formatted report
	report := git.FormatStatsReport(stats)
	fmt.Println(report)

	// Show top patterns
	patterns := git.GetTopPatterns(stats, 3)
	if len(patterns) > 0 {
		fmt.Println("💡 Your Top Commit Patterns:")
		for i, pattern := range patterns {
			fmt.Printf("  %d. %s", i+1, pattern.Type)
			if pattern.Scope != "" {
				fmt.Printf("(%s)", pattern.Scope)
			}
			fmt.Printf(" - used %d times\n", pattern.Frequency)
		}
		fmt.Println()
	}

	// Export to JSON if requested
	if statsExport != "" {
		if err := exportStatsToJSON(stats, statsExport); err != nil {
			return fmt.Errorf("failed to export stats: %w", err)
		}
		fmt.Printf("✅ Statistics exported to %s\n", statsExport)
	}

	// Show insights and recommendations
	showInsights(stats)

	return nil
}

func showInsights(stats *git.CommitStats) {
	fmt.Println("💭 Insights & Recommendations:")

	// Check subject length
	if stats.AverageLength > 72 {
		fmt.Println("  ⚠️  Your average subject line is quite long (>72 chars)")
		fmt.Println("     Consider using shorter, more concise subjects")
	} else if stats.AverageLength < 30 {
		fmt.Println("  ℹ️  Your subject lines are very brief (<30 chars)")
		fmt.Println("     Consider adding more context when helpful")
	}

	// Check scope usage
	scopePercentage := float64(stats.WithScope) / float64(stats.TotalCommits) * 100
	if scopePercentage < 20 {
		fmt.Println("  💡 You rarely use scopes in commits (<20%)")
		fmt.Println("     Scopes help organize changes by component/module")
	}

	// Check body usage
	bodyPercentage := float64(stats.WithBody) / float64(stats.TotalCommits) * 100
	if bodyPercentage < 10 {
		fmt.Println("  💡 Most commits have no body (<10%)")
		fmt.Println("     Consider adding details for non-trivial changes")
	}

	// Check commit frequency
	if stats.RecentTrends != nil {
		if stats.RecentTrends.AveragePerDay > 10 {
			fmt.Println("  🔥 Very high commit frequency (>10/day avg)")
			fmt.Println("     Great activity! Consider squashing related commits")
		} else if stats.RecentTrends.AveragePerDay < 1 {
			fmt.Println("  📉 Low commit frequency (<1/day avg)")
			fmt.Println("     Consider committing more frequently")
		}
	}

	// Check type diversity
	if len(stats.TypeDistribution) < 3 {
		fmt.Println("  ℹ️  Limited commit type variety")
		fmt.Println("     Explore other types: docs, test, refactor, perf, etc.")
	}

	fmt.Println()
}

func exportStatsToJSON(stats *git.CommitStats, filename string) error {
	// This is a placeholder for JSON export functionality
	// In a real implementation, you would use encoding/json
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Simple JSON export (in practice, use json.Marshal)
	_, err = file.WriteString(fmt.Sprintf(`{
  "total_commits": %d,
  "average_length": %d,
  "with_scope": %d,
  "with_body": %d,
  "with_ticket": %d,
  "type_distribution": %v,
  "scope_distribution": %v,
  "language_usage": %v
}`, stats.TotalCommits, stats.AverageLength, stats.WithScope, stats.WithBody,
		stats.WithTicket, stats.TypeDistribution, stats.ScopeDistribution, stats.LanguageUsage))

	return err
}
