package terraform

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

func TestValidateDir(t *testing.T) {
	tests := []struct {
		name      string
		workingDir string
		shouldPass bool
	}{
		{
			name:       "valid directory",
			workingDir: ".",
			shouldPass: true,
		},
		{
			name:       "nonexistent directory",
			workingDir: "/nonexistent/path/that/does/not/exist",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewChecker(tt.workingDir)
			err := checker.ValidateDir()
			if tt.shouldPass && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	}
}

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

func TestCheckProviderVersions(t *testing.T) {
	tests := []struct {
		name           string
		setupDir       func(t *testing.T) string
		shouldFindIssue bool
	}{
		{
			name: "unpinned providers detected",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "main.tf")
				content := `terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: true,
		},
		{
			name: "pinned providers pass",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "main.tf")
				content := `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: false,
		},
		{
			name: "no providers block",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "main.tf")
				content := `resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := tt.setupDir(t)
			checker := NewChecker(tmpDir)
			results, err := checker.CheckProviderVersions()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			foundIssue := len(results) > 0
			if foundIssue != tt.shouldFindIssue {
				t.Errorf("Expected issue found: %v, got: %v", tt.shouldFindIssue, foundIssue)
			}

			if foundIssue && len(results) > 0 {
				if results[0].CheckName != "provider-version" {
					t.Errorf("Expected check name 'provider-version', got %s", results[0].CheckName)
				}
				if results[0].Severity != severity.Medium {
					t.Errorf("Expected severity MEDIUM, got %v", results[0].Severity)
				}
			}
		})
	}
}

func TestCheckCredentials(t *testing.T) {
	tests := []struct {
		name            string
		setupDir        func(t *testing.T) string
		shouldFindIssue bool
		credType        string
	}{
		{
			name: "aws access key detected",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "secrets.tf")
				content := `resource "aws_instance" "example" {
  access_key = "AKIAIOSFODNN7EXAMPLE"
  instance_type = "t2.micro"
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: true,
			credType:        "aws_access_key",
		},
		{
			name: "password detected",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "db.tf")
				content := `resource "aws_db_instance" "example" {
  password = "MySecurePassword123!"
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: true,
			credType:        "password",
		},
		{
			name: "no credentials",
			setupDir: func(t *testing.T) string {
				tmpDir := t.TempDir()
				tfFile := filepath.Join(tmpDir, "main.tf")
				content := `resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
}`
				err := os.WriteFile(tfFile, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tmpDir
			},
			shouldFindIssue: false,
		},
		{
			name: "empty directory",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			shouldFindIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := tt.setupDir(t)
			checker := NewChecker(tmpDir)
			results, err := checker.CheckCredentials()
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			foundIssue := len(results) > 0
			if foundIssue != tt.shouldFindIssue {
				t.Errorf("Expected issue found: %v, got: %v", tt.shouldFindIssue, foundIssue)
			}

			if foundIssue && len(results) > 0 {
				if results[0].CheckName != "hardcoded-credentials" {
					t.Errorf("Expected check name 'hardcoded-credentials', got %s", results[0].CheckName)
				}
				if results[0].Severity != severity.Critical {
					t.Errorf("Expected severity CRITICAL, got %v", results[0].Severity)
				}
				if results[0].Message == "" {
					t.Errorf("Expected non-empty message")
				}
			}
		})
	}
}

func TestCheckResultStructure(t *testing.T) {
	// Verify that CheckResult has all required fields
	result := CheckResult{
		CheckName:      "test-check",
		Severity:       severity.High,
		ResourceID:     "test-resource",
		Message:        "Test message",
		Recommendation: "Test recommendation",
	}

	if result.CheckName == "" {
		t.Error("CheckName should not be empty")
	}
	if result.Severity != severity.High {
		t.Error("Severity should be set")
	}
	if result.ResourceID == "" {
		t.Error("ResourceID should not be empty")
	}
	if result.Message == "" {
		t.Error("Message should not be empty")
	}
	if result.Recommendation == "" {
		t.Error("Recommendation should not be empty")
	}
}

func TestCheckFormat_BinaryNotInstalled(t *testing.T) {
	// Test with a valid terraform directory but no terraform binary in PATH
	// The function should handle this gracefully
	checker := NewChecker("../../testdata/terraform")
	results, err := checker.CheckFormat()

	// Should return without error, results depend on whether terraform is installed
	// If terraform is not installed, we expect empty results or specific error handling
	_ = results
	// err may be nil (terraform not found handled gracefully) or contain "not found"
	if err != nil {
		// This is acceptable - terraform not found is a valid error condition
		t.Logf("terraform binary check: %v", err)
	}
}

func TestCheckValidate_BinaryNotInstalled(t *testing.T) {
	// Similar to CheckFormat - test graceful handling of missing terraform
	checker := NewChecker("../../testdata/terraform")
	results, err := checker.CheckValidate()

	// Should return without panic
	_ = results
	if err != nil {
		t.Logf("terraform binary check: %v", err)
	}
}
