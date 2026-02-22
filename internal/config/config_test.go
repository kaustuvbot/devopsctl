package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.AWS.Region != "us-east-1" {
		t.Errorf("expected default region us-east-1, got %s", cfg.AWS.Region)
	}
	if cfg.AWS.KeyAgeDays != 90 {
		t.Errorf("expected default key age 90, got %d", cfg.AWS.KeyAgeDays)
	}
	if cfg.Git.RepoSizeMB != 500 {
		t.Errorf("expected default repo size 500, got %d", cfg.Git.RepoSizeMB)
	}
	if cfg.Git.BranchAgeDays != 90 {
		t.Errorf("expected default branch age 90, got %d", cfg.Git.BranchAgeDays)
	}
	if cfg.Git.LargeFileMB != 50 {
		t.Errorf("expected default large file 50, got %d", cfg.Git.LargeFileMB)
	}
	if cfg.Docker.DockerfilePath != "Dockerfile" {
		t.Errorf("expected default dockerfile path Dockerfile, got %s", cfg.Docker.DockerfilePath)
	}
}

func TestLoadValidConfig(t *testing.T) {
	cfg, err := Load(filepath.Join("testdata", "valid.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AWS.Region != "eu-west-1" {
		t.Errorf("expected region eu-west-1, got %s", cfg.AWS.Region)
	}
	if cfg.AWS.KeyAgeDays != 60 {
		t.Errorf("expected key age 60, got %d", cfg.AWS.KeyAgeDays)
	}
	if cfg.AWS.Profile != "staging" {
		t.Errorf("expected profile staging, got %s", cfg.AWS.Profile)
	}
	if cfg.Git.RepoSizeMB != 200 {
		t.Errorf("expected repo size 200, got %d", cfg.Git.RepoSizeMB)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load("nonexistent.yaml")
	if err != nil {
		t.Fatalf("missing file should return defaults, got error: %v", err)
	}

	if cfg.AWS.Region != "us-east-1" {
		t.Errorf("expected default region, got %s", cfg.AWS.Region)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	_, err := Load(filepath.Join("testdata", "invalid.yaml"))
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestFindConfigFile(t *testing.T) {
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Logf("failed to restore dir: %v", err)
		}
	}()

	// No config file exists
	result := FindConfigFile()
	if result != "" {
		t.Errorf("expected empty string when no config exists, got %s", result)
	}

	// Create .devopsctl.yaml
	if err := os.WriteFile(".devopsctl.yaml", []byte("aws:\n  region: test\n"), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}
	result = FindConfigFile()
	if result != ".devopsctl.yaml" {
		t.Errorf("expected .devopsctl.yaml, got %s", result)
	}
}

func TestDefaultConfig_EnabledFlags(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.AWS.Enabled {
		t.Error("expected AWS.Enabled to default to true")
	}
	if !cfg.Docker.Enabled {
		t.Error("expected Docker.Enabled to default to true")
	}
	if !cfg.Terraform.Enabled {
		t.Error("expected Terraform.Enabled to default to true")
	}
	if !cfg.Git.Enabled {
		t.Error("expected Git.Enabled to default to true")
	}
}

func TestDefaultConfig_IgnoreChecks(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Ignore.Checks == nil {
		t.Error("expected Ignore.Checks to be empty slice, got nil")
	}
	if len(cfg.Ignore.Checks) != 0 {
		t.Errorf("expected Ignore.Checks to have length 0, got %d", len(cfg.Ignore.Checks))
	}
}

func TestLoad_PartialOverride(t *testing.T) {
	// Create a temp config file with partial overrides
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".devopsctl.yaml")
	configContent := `
aws:
  enabled: false
docker:
  enabled: true
ignore:
  checks:
    - iam-mfa-disabled
    - s3-public-bucket
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// AWS should be disabled
	if cfg.AWS.Enabled != false {
		t.Errorf("expected AWS.Enabled to be false, got %v", cfg.AWS.Enabled)
	}

	// Docker should still be enabled (not set in yaml, should keep default)
	if !cfg.Docker.Enabled {
		t.Error("expected Docker.Enabled to be true (default)")
	}

	// Terraform not set, should keep default
	if !cfg.Terraform.Enabled {
		t.Error("expected Terraform.Enabled to be true (default)")
	}

	// Git not set, should keep default
	if !cfg.Git.Enabled {
		t.Error("expected Git.Enabled to be true (default)")
	}

	// Ignore checks should be populated
	if len(cfg.Ignore.Checks) != 2 {
		t.Errorf("expected 2 ignore checks, got %d", len(cfg.Ignore.Checks))
	}
	if cfg.Ignore.Checks[0] != "iam-mfa-disabled" {
		t.Errorf("expected first ignore check to be iam-mfa-disabled, got %s", cfg.Ignore.Checks[0])
	}
	if cfg.Ignore.Checks[1] != "s3-public-bucket" {
		t.Errorf("expected second ignore check to be s3-public-bucket, got %s", cfg.Ignore.Checks[1])
	}
}
