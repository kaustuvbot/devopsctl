package terraform

import (
	"os/exec"

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
