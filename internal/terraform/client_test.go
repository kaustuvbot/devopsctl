package terraform

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	workingDir := "/path/to/terraform"
	client := NewClient(workingDir)

	if client == nil {
		t.Errorf("Expected non-nil client")
	}
	if client.workingDir != workingDir {
		t.Errorf("Expected workingDir to be %s, got %s", workingDir, client.workingDir)
	}
}