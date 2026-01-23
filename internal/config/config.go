package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// AWSConfig holds AWS-specific configuration.
type AWSConfig struct {
	Enabled      bool   `yaml:"enabled"`
	Region       string `yaml:"region"`
	Profile      string `yaml:"profile"`
	KeyAgeDays   int    `yaml:"key_age_days"`
}

// DockerConfig holds Docker-specific configuration.
type DockerConfig struct {
	Enabled        bool   `yaml:"enabled"`
	DockerfilePath string `yaml:"dockerfile_path"`
}

// TerraformConfig holds Terraform-specific configuration.
type TerraformConfig struct {
	Enabled bool   `yaml:"enabled"`
	TfDir   string `yaml:"tf_dir"`
}

// GitConfig holds Git-specific configuration.
type GitConfig struct {
	Enabled       bool `yaml:"enabled"`
	RepoSizeMB    int  `yaml:"repo_size_mb"`
	BranchAgeDays int  `yaml:"branch_age_days"`
	LargeFileMB   int  `yaml:"large_file_mb"`
}

// IgnoreConfig holds ignore patterns for check filtering.
type IgnoreConfig struct {
	Checks []string `yaml:"checks"`
}

// Config represents the main configuration structure.
type Config struct {
	AWS       AWSConfig       `yaml:"aws"`
	Docker    DockerConfig    `yaml:"docker"`
	Terraform TerraformConfig `yaml:"terraform"`
	Git       GitConfig       `yaml:"git"`
	Ignore    IgnoreConfig    `yaml:"ignore"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		AWS: AWSConfig{
			Enabled:    true,
			Region:     "us-east-1",
			KeyAgeDays: 90,
		},
		Docker: DockerConfig{
			Enabled:        true,
			DockerfilePath: "Dockerfile",
		},
		Terraform: TerraformConfig{
			Enabled: true,
		},
		Git: GitConfig{
			Enabled:       true,
			RepoSizeMB:    500,
			BranchAgeDays: 90,
			LargeFileMB:   50,
		},
		Ignore: IgnoreConfig{
			Checks: []string{},
		},
	}
}

// Load reads and parses a YAML config file. If the file does not exist,
// default values are returned.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// FindConfigFile looks for .devopsctl.yaml in standard locations.
func FindConfigFile() string {
	candidates := []string{
		".devopsctl.yaml",
		".devopsctl.yml",
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	return ""
}
