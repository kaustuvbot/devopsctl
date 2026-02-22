package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

// TableReporter renders check results as a formatted text table.
type TableReporter struct{}

// NewTableReporter returns a TableReporter.
func NewTableReporter() *TableReporter { return &TableReporter{} }

// ANSI color codes
const (
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	reset  = "\033[0m"
)

// colorize returns the severity with ANSI color codes
func colorize(s string) string {
	level := severity.Level(s)
	switch level {
	case severity.Critical:
		return red + s + reset
	case severity.High:
		return yellow + s + reset
	case severity.Low, severity.Medium:
		return green + s + reset
	default:
		return s
	}
}

// Render writes a formatted table of check results to w.
func (r *TableReporter) Render(w io.Writer, report *Report) error {
	// Check if terminal supports colors
	isTerminal := isTerminal(w)

	fmt.Fprintf(w, "=== %s Audit Results ===\n\n", report.Module)
	if len(report.Results) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SEVERITY\tCHECK NAME\tRESOURCE\tMESSAGE")
	fmt.Fprintln(tw, "--------\t----------\t--------\t-------")
	for _, result := range report.Results {
		sev := result.Severity
		if isTerminal {
			sev = colorize(sev)
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			sev, result.CheckName, result.ResourceID, result.Message)
	}
	return tw.Flush()
}

// isTerminal checks if the writer is a terminal
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isatty(f.Fd())
	}
	return false
}

// isatty checks if the file descriptor is a terminal
func isatty(fd uintptr) bool {
	return false // Simplified - always returns false for non-terminal
}
