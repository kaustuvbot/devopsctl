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
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	// No config file exists
	result := FindConfigFile()
	if result != "" {
		t.Errorf("expected empty string when no config exists, got %s", result)
	}

	// Create .devopsctl.yaml
	os.WriteFile(".devopsctl.yaml", []byte("aws:\n  region: test\n"), 0644)
	result = FindConfigFile()
	if result != ".devopsctl.yaml" {
		t.Errorf("expected .devopsctl.yaml, got %s", result)
	}
}
