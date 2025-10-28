package reporter

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableReporter_EmptyResults(t *testing.T) {
	rep := NewTableReporter()
	var buf bytes.Buffer
	report := &Report{Module: "aws", Results: []CheckResult{}}
	if err := rep.Render(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No issues found") {
		t.Errorf("expected 'No issues found' message, got: %s", buf.String())
	}
}

func TestTableReporter_WithResults(t *testing.T) {
	rep := NewTableReporter()
	var buf bytes.Buffer
	report := &Report{
		Module: "aws",
		Results: []CheckResult{
			{CheckName: "iam-mfa-disabled", Severity: "HIGH", ResourceID: "alice", Message: "No MFA"},
			{CheckName: "s3-public-bucket", Severity: "CRITICAL", ResourceID: "my-bucket", Message: "Public"},
		},
	}
	if err := rep.Render(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"HIGH", "CRITICAL", "alice", "my-bucket", "iam-mfa-disabled"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestTableReporter_HeaderPresent(t *testing.T) {
	rep := NewTableReporter()
	var buf bytes.Buffer
	report := &Report{Module: "docker", Results: []CheckResult{
		{CheckName: "test", Severity: "LOW", ResourceID: "res", Message: "msg"},
	}}
	if err := rep.Render(&buf, report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "docker Audit Results") {
		t.Errorf("expected module name in header, got: %s", out)
	}
	if !strings.Contains(out, "SEVERITY") {
		t.Errorf("expected SEVERITY column header, got: %s", out)
	}
}
