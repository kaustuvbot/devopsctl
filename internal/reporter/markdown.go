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
	if _, err := fmt.Fprintf(w, "# %s Audit Report\n\n", titleCase(report.Module)); err != nil {
		return err
	}

	if len(report.Results) == 0 {
		_, err := fmt.Fprintf(w, "No findings.\n\n")
		return err
	}

	// Create a tabwriter for alignment
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(tw, "| %s | %s | %s | %s |\n", "Severity", "Check", "Resource", "Message")
	_, _ = fmt.Fprintf(tw, "| --- | --- | --- | --- |\n")

	for _, result := range report.Results {
		_, _ = fmt.Fprintf(tw, "| %s | %s | %s | %s |\n",
			result.Severity,
			result.CheckName,
			result.ResourceID,
			result.Message,
		)
	}

	_ = tw.Flush()
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}

	// Add recommendations section
	if _, err := fmt.Fprintf(w, "## Recommendations\n\n"); err != nil {
		return err
	}
	for _, result := range report.Results {
		if result.Recommendation != "" {
			if _, err := fmt.Fprintf(w, "- **%s**: %s\n", result.CheckName, result.Recommendation); err != nil {
				return err
			}
		}
	}

	_, err := fmt.Fprintf(w, "\n")
	return err
}

// titleCase converts a string to title case.
func titleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
