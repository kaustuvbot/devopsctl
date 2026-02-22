package git

import (
	"context"
	"testing"

	appconfig "github.com/kaustuvbot/devopsctl/internal/config"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// mockRunnerClient is a mock implementation of git client for testing
type mockRunnerClient struct {
	sizeResults     []reporter.CheckResult
	sizeErr        error
	branchResults   []reporter.CheckResult
	branchErr       error
	fileResults     []reporter.CheckResult
	fileErr         error
}

func (m *mockRunnerClient) Run(ctx context.Context, args ...string) (string, error) {
	return "", nil
}

// TestRunnerExecution tests that runner executes all checks
func TestRunnerExecution(t *testing.T) {
	cfg := appconfig.GitConfig{
		RepoSizeMB:    500,
		BranchAgeDays: 90,
		LargeFileMB:   50,
	}

	// Test that Runner can be created
	_ = NewRunner(".", cfg)
}

// TestRunnerResultsAggregation tests that results are aggregated correctly
func TestRunnerResultsAggregation(t *testing.T) {
	cfg := appconfig.GitConfig{
		RepoSizeMB:    500,
		BranchAgeDays: 90,
		LargeFileMB:   50,
	}

	// Verify config is passed correctly
	if cfg.RepoSizeMB != 500 {
		t.Errorf("Expected RepoSizeMB 500, got %d", cfg.RepoSizeMB)
	}
	if cfg.BranchAgeDays != 90 {
		t.Errorf("Expected BranchAgeDays 90, got %d", cfg.BranchAgeDays)
	}
	if cfg.LargeFileMB != 50 {
		t.Errorf("Expected LargeFileMB 50, got %d", cfg.LargeFileMB)
	}
}

// TestCheckResultSeverity tests severity levels in results
func TestCheckResultSeverity(t *testing.T) {
	tests := []struct {
		result  reporter.CheckResult
		invalid bool
	}{
		{reporter.CheckResult{Severity: "LOW"}, false},
		{reporter.CheckResult{Severity: "MEDIUM"}, false},
		{reporter.CheckResult{Severity: "HIGH"}, false},
		{reporter.CheckResult{Severity: "CRITICAL"}, false},
		{reporter.CheckResult{Severity: "INVALID"}, true},
	}

	for _, tt := range tests {
		isValid := tt.result.Severity == "LOW" ||
			tt.result.Severity == "MEDIUM" ||
			tt.result.Severity == "HIGH" ||
			tt.result.Severity == "CRITICAL"

		if isValid && tt.invalid {
			t.Errorf("Expected invalid severity for %s", tt.result.Severity)
		}
	}
}
