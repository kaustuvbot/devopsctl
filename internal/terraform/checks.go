package terraform

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kaustuvbot/devopsctl/internal/severity"
)

// CheckResult represents the result of a terraform check.
type CheckResult struct {
	CheckName      string
	Severity       severity.Level
	ResourceID     string
	Message        string
	Recommendation string
}

// Checker performs terraform validation checks.
type Checker struct {
	workingDir string
}

// NewChecker creates a new terraform checker.
func NewChecker(workingDir string) *Checker {
	return &Checker{workingDir: workingDir}
}

// ValidateDir checks if the working directory exists and is readable.
func (c *Checker) ValidateDir() error {
	info, err := os.Stat(c.workingDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return os.ErrNotExist
	}
	return nil
}

// CheckFormat runs terraform fmt -check to verify formatting.
func (c *Checker) CheckFormat() ([]CheckResult, error) {
	// Check if terraform binary exists
	_, err := exec.LookPath("terraform")
	if err != nil {
		// terraform not found, skip this check
		return []CheckResult{}, nil
	}

	cmd := exec.Command("terraform", "fmt", "-check", "-recursive")
	cmd.Dir = c.workingDir
	err = cmd.Run()

	var results []CheckResult
	if err != nil {
		results = append(results, CheckResult{
			CheckName:      "terraform-fmt",
			Severity:       severity.Medium,
			ResourceID:     c.workingDir,
			Message:        "Terraform files are not properly formatted",
			Recommendation: "Run 'terraform fmt' to fix formatting",
		})
	}
	return results, nil
}

// CheckValidate runs terraform validate to check configuration validity.
func (c *Checker) CheckValidate() ([]CheckResult, error) {
	// Check if terraform binary exists
	_, err := exec.LookPath("terraform")
	if err != nil {
		// terraform not found, skip this check
		return []CheckResult{}, nil
	}

	cmd := exec.Command("terraform", "validate")
	cmd.Dir = c.workingDir
	err = cmd.Run()

	var results []CheckResult
	if err != nil {
		results = append(results, CheckResult{
			CheckName:      "terraform-validate",
			Severity:       severity.High,
			ResourceID:     c.workingDir,
			Message:        "Terraform configuration is invalid",
			Recommendation: "Fix terraform validation errors",
		})
	}
	return results, nil
}

// CheckProviderVersions checks for unpinned provider versions in terraform files.
func (c *Checker) CheckProviderVersions() ([]CheckResult, error) {
	var results []CheckResult

	files, err := filepath.Glob(filepath.Join(c.workingDir, "*.tf"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Simple check: if file has required_providers but no version constraint
		hasRequiredProviders := false
		hasVersion := false

		contentStr := string(content)
		if strings.Contains(contentStr, "required_providers") {
			hasRequiredProviders = true
		}
		if strings.Contains(contentStr, "version") {
			hasVersion = true
		}

		if hasRequiredProviders && !hasVersion {
			results = append(results, CheckResult{
				CheckName:      "provider-version",
				Severity:       severity.Medium,
				ResourceID:     file,
				Message:        "Provider version constraint not found",
				Recommendation: "Add version constraint to provider configuration",
			})
		}
	}

	return results, nil
}

// CheckCredentials detects hardcoded credentials in terraform files.
func (c *Checker) CheckCredentials() ([]CheckResult, error) {
	var results []CheckResult

	// Patterns for detecting hardcoded credentials
	patterns := map[string]string{
		"aws_access_key": `AKIA[0-9A-Z]{16}`,
		"aws_secret_key": `[A-Za-z0-9/+=]{40}`,
		"password":       `password\s*=\s*"[^"]+"`,
		"api_key":        `api_key\s*=\s*"[^"]+"`,
		"secret":         `secret\s*=\s*"[^"]+"`,
	}

	files, err := filepath.Glob(filepath.Join(c.workingDir, "*.tf"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		contentStr := string(content)

		for credType, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(contentStr) {
				results = append(results, CheckResult{
					CheckName:      "hardcoded-credentials",
					Severity:       severity.Critical,
					ResourceID:     file,
					Message:        "Hardcoded " + credType + " detected",
					Recommendation: "Use environment variables or secret management instead",
				})
			}
		}
	}

	return results, nil
}
