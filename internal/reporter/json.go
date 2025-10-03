package reporter

import (
	"encoding/json"
	"io"
)

// JSONReporter outputs results in JSON format.
type JSONReporter struct {
	Pretty bool
}

// NewJSONReporter creates a new JSON reporter.
func NewJSONReporter(pretty bool) *JSONReporter {
	return &JSONReporter{Pretty: pretty}
}

// Render writes the report as JSON to the given writer.
func (r *JSONReporter) Render(w io.Writer, report *Report) error {
	encoder := json.NewEncoder(w)
	if r.Pretty {
		encoder.SetIndent("", "  ")
	}
	return encoder.Encode(report)
}
