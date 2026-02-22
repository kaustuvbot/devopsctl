package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONReporter_Render(t *testing.T) {
	tests := []struct {
		name     string
		pretty   bool
		module   string
		results  []CheckResult
		wantKeys []string // top-level keys we expect in the JSON
	}{
		{
			name:   "empty results",
			pretty: true,
			module: "aws",
			results: []CheckResult{},
			wantKeys: []string{"module", "results"},
		},
		{
			name:   "single result with all fields",
			pretty: true,
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
			wantKeys: []string{"module", "results"},
		},
		{
			name:   "multiple results",
			pretty: true,
			module: "git",
			results: []CheckResult{
				{CheckName: "git-stale-branch", Severity: "LOW", ResourceID: "old-branch"},
				{CheckName: "git-repo-size", Severity: "MEDIUM", ResourceID: "."},
			},
			wantKeys: []string{"module", "results"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := NewJSONReporter(tt.pretty)
			report := &Report{Module: tt.module, Results: tt.results}

			var buf bytes.Buffer
			err := reporter.Render(&buf, report)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()

			// Verify JSON is valid
			var parsed map[string]interface{}
			err = json.Unmarshal([]byte(output), &parsed)
			if err != nil {
				t.Fatalf("invalid JSON output: %v", err)
			}

			// Verify expected top-level keys exist
			for _, key := range tt.wantKeys {
				if _, ok := parsed[key]; !ok {
					t.Errorf("expected key %q in output", key)
				}
			}

			// Verify "module" value
			if parsed["module"] != tt.module {
				t.Errorf("expected module %q, got %v", tt.module, parsed["module"])
			}
		})
	}
}

func TestJSONReporter_PrettyVsCompact(t *testing.T) {
	report := &Report{
		Module: "test",
		Results: []CheckResult{
			{CheckName: "test-check", Severity: "HIGH", ResourceID: "res1"},
		},
	}

	// Pretty mode
	prettyReporter := NewJSONReporter(true)
	var prettyBuf bytes.Buffer
	if err := prettyReporter.Render(&prettyBuf, report); err != nil {
		t.Fatalf("pretty render failed: %v", err)
	}
	prettyOutput := prettyBuf.String()

	// Compact mode
	compactReporter := NewJSONReporter(false)
	var compactBuf bytes.Buffer
	if err := compactReporter.Render(&compactBuf, report); err != nil {
		t.Fatalf("compact render failed: %v", err)
	}
	compactOutput := compactBuf.String()

	// Pretty output should have indentation
	if !strings.Contains(prettyOutput, "  ") {
		t.Error("expected pretty output to contain indentation")
	}

	// Compact output should not have newlines except in values
	lines := strings.Split(compactOutput, "\n")
	if len(lines) > 2 {
		t.Errorf("expected compact output to be on single line, got %d lines", len(lines))
	}
}

func TestJSONReporter_CheckResultFields(t *testing.T) {
	report := &Report{
		Module: "aws",
		Results: []CheckResult{
			{
				CheckName:      "check-name",
				Severity:       "CRITICAL",
				ResourceID:     "resource-id",
				Message:        "test message",
				Recommendation: "test recommendation",
			},
		},
	}

	reporter := NewJSONReporter(true)
	var buf bytes.Buffer
	err := reporter.Render(&buf, report)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify individual fields are present in the JSON
	output := buf.String()
	if !strings.Contains(output, "check_name") {
		t.Error("expected 'check_name' in JSON output")
	}
	if !strings.Contains(output, "severity") {
		t.Error("expected 'severity' in JSON output")
	}
	if !strings.Contains(output, "resource_id") {
		t.Error("expected 'resource_id' in JSON output")
	}
	if !strings.Contains(output, "message") {
		t.Error("expected 'message' in JSON output")
	}
	if !strings.Contains(output, "recommendation") {
		t.Error("expected 'recommendation' in JSON output")
	}
}
