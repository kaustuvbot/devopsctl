package reporter

import "io"

// CheckResult represents the output of a single check.
type CheckResult struct {
	CheckName      string `json:"check_name"`
	Severity       string `json:"severity"`
	ResourceID     string `json:"resource_id"`
	Message        string `json:"message"`
	Recommendation string `json:"recommendation"`
}

// Report holds a collection of check results for a module.
type Report struct {
	Module  string        `json:"module"`
	Results []CheckResult `json:"results"`
}

// Reporter defines the interface for output formatting.
type Reporter interface {
	Render(w io.Writer, report *Report) error
}
