package terraform

import "fmt"

// Runner orchestrates terraform validation checks.
type Runner struct {
	workingDir string
	checker    *Checker
}

// NewRunner creates a new terraform runner.
func NewRunner(workingDir string) *Runner {
	return &Runner{
		workingDir: workingDir,
		checker:    NewChecker(workingDir),
	}
}

// RunAllChecks runs all terraform checks and returns combined results.
// Checks that fail are skipped, not fatal - returns partial results.
func (r *Runner) RunAllChecks() ([]CheckResult, error) {
	var allResults []CheckResult
	var errs []string

	// Run format check
	formatResults, err := r.checker.CheckFormat()
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		allResults = append(allResults, formatResults...)
	}

	// Run validate check
	validateResults, err := r.checker.CheckValidate()
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		allResults = append(allResults, validateResults...)
	}

	// Run provider version check
	providerResults, err := r.checker.CheckProviderVersions()
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		allResults = append(allResults, providerResults...)
	}

	// Run credentials check
	credResults, err := r.checker.CheckCredentials()
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		allResults = append(allResults, credResults...)
	}

	if len(errs) > 0 {
		return allResults, fmt.Errorf("some checks failed: %v", errs)
	}
	return allResults, nil
}