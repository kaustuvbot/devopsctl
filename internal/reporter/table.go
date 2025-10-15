package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// TableReporter renders check results as a formatted text table.
type TableReporter struct{}

// NewTableReporter returns a TableReporter.
func NewTableReporter() *TableReporter { return &TableReporter{} }

// Render writes a formatted table of check results to w.
func (r *TableReporter) Render(w io.Writer, report *Report) error {
	fmt.Fprintf(w, "=== %s Audit Results ===\n\n", report.Module)
	if len(report.Results) == 0 {
		fmt.Fprintln(w, "No issues found.")
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SEVERITY\tCHECK NAME\tRESOURCE\tMESSAGE")
	fmt.Fprintln(tw, "--------\t----------\t--------\t-------")
	for _, result := range report.Results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			result.Severity, result.CheckName, result.ResourceID, result.Message)
	}
	return tw.Flush()
}
