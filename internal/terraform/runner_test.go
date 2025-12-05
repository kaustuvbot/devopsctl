package terraform

import (
	"testing"
)

func TestTerraformRunner(t *testing.T) {
	runner := NewRunner("../../testdata/terraform")
	results, err := runner.RunAllChecks()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Results will vary based on terraform installation and file contents
	_ = results
}