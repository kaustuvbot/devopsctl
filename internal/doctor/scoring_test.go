package doctor

import (
	"testing"

	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

func TestComputeSummary_AllSeverities(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "CRITICAL"},
				{CheckName: "check2", Severity: "HIGH"},
				{CheckName: "check3", Severity: "MEDIUM"},
				{CheckName: "check4", Severity: "LOW"},
			},
		},
	}

	summary := ComputeSummary(reports)

	if summary.TotalFindings != 4 {
		t.Errorf("expected 4 total findings, got %d", summary.TotalFindings)
	}
	if summary.Critical != 1 {
		t.Errorf("expected 1 critical, got %d", summary.Critical)
	}
	if summary.High != 1 {
		t.Errorf("expected 1 high, got %d", summary.High)
	}
	if summary.Medium != 1 {
		t.Errorf("expected 1 medium, got %d", summary.Medium)
	}
	if summary.Low != 1 {
		t.Errorf("expected 1 low, got %d", summary.Low)
	}
	// Score = 4 (Critical) + 3 (High) + 2 (Medium) + 1 (Low) = 10
	if summary.Score != 10 {
		t.Errorf("expected score 10, got %d", summary.Score)
	}
	if summary.ModulesFailed != 0 {
		t.Errorf("expected 0 modules failed, got %d", summary.ModulesFailed)
	}
}

func TestComputeSummary_FailedModules(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{},
			Error:   "access denied",
		},
		{
			Module: "docker",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "LOW"},
			},
		},
	}

	summary := ComputeSummary(reports)

	if summary.ModulesFailed != 1 {
		t.Errorf("expected 1 module failed, got %d", summary.ModulesFailed)
	}
	if summary.TotalFindings != 1 {
		t.Errorf("expected 1 total finding, got %d", summary.TotalFindings)
	}
	if summary.ModuleErrors["aws"] != "access denied" {
		t.Errorf("expected 'access denied' error for aws, got %s", summary.ModuleErrors["aws"])
	}
}

func TestComputeSummary_Empty(t *testing.T) {
	reports := []ModuleReport{}

	summary := ComputeSummary(reports)

	if summary.TotalFindings != 0 {
		t.Errorf("expected 0 findings, got %d", summary.TotalFindings)
	}
	if summary.Score != 0 {
		t.Errorf("expected score 0, got %d", summary.Score)
	}
	if summary.ModulesFailed != 0 {
		t.Errorf("expected 0 modules failed, got %d", summary.ModulesFailed)
	}
}

func TestExitCode_Critical(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "CRITICAL"},
			},
		},
	}

	code := ExitCode(reports)
	if code != 4 {
		t.Errorf("expected exit code 4 for CRITICAL, got %d", code)
	}
}

func TestExitCode_High(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "HIGH"},
			},
		},
	}

	code := ExitCode(reports)
	if code != 3 {
		t.Errorf("expected exit code 3 for HIGH, got %d", code)
	}
}

func TestExitCode_Medium(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "MEDIUM"},
			},
		},
	}

	code := ExitCode(reports)
	if code != 2 {
		t.Errorf("expected exit code 2 for MEDIUM, got %d", code)
	}
}

func TestExitCode_Low(t *testing.T) {
	reports := []ModuleReport{
		{
			Module: "aws",
			Results: []reporter.CheckResult{
				{CheckName: "check1", Severity: "LOW"},
			},
		},
	}

	code := ExitCode(reports)
	if code != 1 {
		t.Errorf("expected exit code 1 for LOW, got %d", code)
	}
}

func TestExitCode_NoFindings(t *testing.T) {
	reports := []ModuleReport{}

	code := ExitCode(reports)
	if code != 0 {
		t.Errorf("expected exit code 0 for no findings, got %d", code)
	}
}

func TestHighestSeverity(t *testing.T) {
	tests := []struct {
		name     string
		reports  []ModuleReport
		expected severity.Level
	}{
		{
			name: "critical is highest",
			reports: []ModuleReport{
				{Module: "aws", Results: []reporter.CheckResult{{Severity: "LOW"}}},
				{Module: "docker", Results: []reporter.CheckResult{{Severity: "CRITICAL"}}},
			},
			expected: severity.Critical,
		},
		{
			name: "high beats medium",
			reports: []ModuleReport{
				{Module: "aws", Results: []reporter.CheckResult{{Severity: "MEDIUM"}}},
				{Module: "docker", Results: []reporter.CheckResult{{Severity: "HIGH"}}},
			},
			expected: severity.High,
		},
		{
			name: "empty reports",
			reports: []ModuleReport{
				{Module: "aws", Results: []reporter.CheckResult{}},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HighestSeverity(tt.reports)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
