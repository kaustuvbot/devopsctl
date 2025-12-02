package terraform

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kaustuvprajapati/devopsctl/internal/severity"
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

// CheckFormat runs terraform fmt -check to verify formatting.
func (c *Checker) CheckFormat() ([]CheckResult, error) {
	cmd := exec.Command("terraform", "fmt", "-check", "-recursive")
	cmd.Dir = c.workingDir
	err := cmd.Run()

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
	cmd := exec.Command("terraform", "validate")
	cmd.Dir = c.workingDir
	err := cmd.Run()

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
