package terraform

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	parser := NewParser()

	file, err := parser.ParseFile("../../testdata/terraform/valid.tf")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if file == nil {
		t.Fatalf("Expected non-nil file")
	}
}
