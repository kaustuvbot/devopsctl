package reporter

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// MarkdownReporter renders reports in Markdown format.
type MarkdownReporter struct{}

// NewMarkdownReporter creates a new MarkdownReporter.
func NewMarkdownReporter() *MarkdownReporter {
	return &MarkdownReporter{}
}

// Render outputs the report in Markdown format.
func (r *MarkdownReporter) Render(w io.Writer, report *Report) error {
	fmt.Fprintf(w, "# %s Audit Report\n\n", titleCase(report.Module))

	if len(report.Results) == 0 {
		fmt.Fprintf(w, "No findings.\n\n")
		return nil
	}

	// Create a tabwriter for alignment
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "| %s | %s | %s | %s |\n", "Severity", "Check", "Resource", "Message")
	fmt.Fprintf(tw, "| --- | --- | --- | --- |\n")

	for _, result := range report.Results {
		fmt.Fprintf(tw, "| %s | %s | %s | %s |\n",
			result.Severity,
			result.CheckName,
			result.ResourceID,
			result.Message,
		)
	}

	tw.Flush()
	fmt.Fprintf(w, "\n")

	// Add recommendations section
	fmt.Fprintf(w, "## Recommendations\n\n")
	for _, result := range report.Results {
		if result.Recommendation != "" {
			fmt.Fprintf(w, "- **%s**: %s\n", result.CheckName, result.Recommendation)
		}
	}

	fmt.Fprintf(w, "\n")
	return nil
}

// titleCase converts a string to title case.
func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
