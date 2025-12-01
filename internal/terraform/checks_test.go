package terraform

import (
	"testing"
)

func TestCheckFormat(t *testing.T) {
	checker := NewChecker("../../testdata/terraform")
	results, err := checker.CheckFormat()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Results will vary based on whether terraform is installed and formatting is correct
	// Just verify the function runs without error
	_ = results
}

func TestCheckValidate(t *testing.T) {
	checker := NewChecker("../../testdata/terraform")
	results, err := checker.CheckValidate()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	// Results will vary based on whether terraform is installed and config is valid
	// Just verify the function runs without error
	_ = results
}