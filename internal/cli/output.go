package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
	"github.com/kaustuvprajapati/devopsctl/internal/severity"
	"github.com/spf13/cobra"
)

// resolveWriter returns a writer for command output.
// If --output is set, opens that file; otherwise returns os.Stdout.
func resolveWriter(_ *cobra.Command) (io.WriteCloser, error) {
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("cannot open output file: %w", err)
		}
		return f, nil
	}
	return os.Stdout, nil
}

// resolveReporter returns the appropriate Reporter based on the --format flag.
// Falls back to --json flag for backwards compatibility.
func resolveReporter() reporter.Reporter {
	// Check --format flag first
	switch outputFormat {
	case "json":
		return reporter.NewJSONReporter(true)
	case "markdown":
		return reporter.NewMarkdownReporter()
	case "table":
		return reporter.NewTableReporter()
	}

	// Fallback to deprecated --json flag
	if jsonOutput {
		return reporter.NewJSONReporter(true)
	}
	return reporter.NewTableReporter()
}

// exitCodeForResults returns the highest severity exit code from results.
// Returns 0 if results is empty.
func exitCodeForResults(results []reporter.CheckResult) int {
	var levels []severity.Level
	for _, r := range results {
		levels = append(levels, severity.Level(r.Severity))
	}
	if len(levels) == 0 {
		return 0
	}
	return severity.Highest(levels).ExitCode()
}

// filterByIgnore filters out checks that match any pattern in Ignore.Checks.
func filterByIgnore(results []reporter.CheckResult, ignorePatterns []string) []reporter.CheckResult {
	if len(ignorePatterns) == 0 {
		return results
	}
	var filtered []reporter.CheckResult
	for _, r := range results {
		skip := false
		for _, pattern := range ignorePatterns {
			if r.CheckName == pattern {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// filterBySeverity filters results to only include CRITICAL and HIGH severity
// when quiet mode is enabled.
func filterBySeverity(results []reporter.CheckResult, quiet bool) []reporter.CheckResult {
	if !quiet {
		return results
	}
	var filtered []reporter.CheckResult
	for _, r := range results {
		if r.Severity == string(severity.Critical) || r.Severity == string(severity.High) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
