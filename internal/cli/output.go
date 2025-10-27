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

// resolveReporter returns the appropriate Reporter based on the --json flag.
func resolveReporter() reporter.Reporter {
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
