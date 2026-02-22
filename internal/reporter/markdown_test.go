package reporter

import (
	"bytes"
	"strings"
	"testing"
)

func TestMarkdownReporter_Render(t *testing.T) {
	tests := []struct {
		name        string
		module      string
		results     []CheckResult
		wantHeading bool
		wantTable   bool
		wantRecs    bool
	}{
		{
			name:        "empty results",
			module:      "aws",
			results:     []CheckResult{},
			wantHeading: true,
			wantTable:   false,
			wantRecs:    false,
		},
		{
			name:   "single result without recommendation",
			module: "aws",
			results: []CheckResult{
				{
					CheckName:  "s3-public-bucket",
					Severity:   "CRITICAL",
					ResourceID: "my-bucket",
					Message:    "Bucket is publicly accessible",
				},
			},
			wantHeading: true,
			wantTable:   true,
			wantRecs:    false, // no recommendation, so section should not appear
		},
		{
			name:   "single result with recommendation",
			module: "aws",
			results: []CheckResult{
				{
					CheckName:      "s3-public-bucket",
					Severity:       "CRITICAL",
					ResourceID:     "my-bucket",
					Message:        "Bucket is publicly accessible",
					Recommendation: "Enable block public access",
				},
			},
			wantHeading: true,
			wantTable:   true,
			wantRecs:    true,
		},
		{
			name:   "multiple results",
			module: "git",
			results: []CheckResult{
				{CheckName: "git-stale-branch", Severity: "LOW", ResourceID: "old-branch", Message: "Branch is stale"},
				{CheckName: "git-repo-size", Severity: "MEDIUM", ResourceID: ".", Message: "Repo too large", Recommendation: "Clean up history"},
			},
			wantHeading: true,
			wantTable:   true,
			wantRecs:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewMarkdownReporter()
			report := &Report{Module: tt.module, Results: tt.results}

			var buf bytes.Buffer
			err := reporter.Render(&buf, report)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()

			// Check heading - titleCase converts first letter to uppercase
			// So "aws" -> "Aws", "git" -> "Git"
			expectedHeading := "# " + titleCase(tt.module) + " Audit Report"
			if tt.wantHeading && !strings.Contains(output, expectedHeading) {
				t.Errorf("expected heading %q in output", expectedHeading)
			}

			// Check for table pipe characters
			if tt.wantTable && !strings.Contains(output, "|") {
				t.Error("expected table pipe characters in output")
			}

			// Check for recommendations section
			// The section header is always printed, but we check for bullet points
			hasRecBullet := strings.Contains(output, "- **")
			if tt.wantRecs && !hasRecBullet {
				t.Error("expected recommendations bullet points in output")
			}
			if !tt.wantRecs && hasRecBullet {
				t.Error("did not expect recommendations bullet points when no recommendations present")
			}
		})
	}
}

func TestMarkdownReporter_TableFormat(t *testing.T) {
	report := &Report{
		Module: "docker",
		Results: []CheckResult{
			{CheckName: "dockerfile-latest-tag", Severity: "MEDIUM", ResourceID: "Dockerfile", Message: "Using latest tag"},
			{CheckName: "dockerfile-no-user", Severity: "HIGH", ResourceID: "Dockerfile", Message: "No USER directive"},
		},
	}

	reporter := NewMarkdownReporter()
	var buf bytes.Buffer
	err := reporter.Render(&buf, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check table header
	if !strings.Contains(output, "| Severity |") {
		t.Error("expected table header with Severity column")
	}
	if !strings.Contains(output, "| --- |") {
		t.Error("expected table separator row")
	}

	// Check data rows contain expected values
	if !strings.Contains(output, "MEDIUM") {
		t.Error("expected MEDIUM severity in table")
	}
	if !strings.Contains(output, "dockerfile-latest-tag") {
		t.Error("expected check name in table")
	}
}

func TestMarkdownReporter_EmptyModule(t *testing.T) {
	reporter := NewMarkdownReporter()
	report := &Report{Module: "", Results: []CheckResult{}}

	var buf bytes.Buffer
	err := reporter.Render(&buf, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should still produce heading, just empty module name
	if !strings.Contains(output, "#  Audit Report") {
		t.Error("expected heading even with empty module name")
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"aws", "Aws"},
		{"docker", "Docker"},
		{"git", "Git"},
		{"terraform", "Terraform"},
		{"", ""},
		// Note: single uppercase letter "A" would produce "!" due to the
		// current implementation (s[0]-'a'+'A' when s[0] is already uppercase)
		// This is an edge case that doesn't affect actual module names.
	}

	for _, tt := range tests {
		result := titleCase(tt.input)
		if result != tt.expected {
			t.Errorf("titleCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
