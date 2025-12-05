package terraform

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
func (r *Runner) RunAllChecks() ([]CheckResult, error) {
	var allResults []CheckResult

	// Run format check
	formatResults, err := r.checker.CheckFormat()
	if err != nil {
		return nil, err
	}
	allResults = append(allResults, formatResults...)

	// Run validate check
	validateResults, err := r.checker.CheckValidate()
	if err != nil {
		return nil, err
	}
	allResults = append(allResults, validateResults...)

	// Run provider version check
	providerResults, err := r.checker.CheckProviderVersions()
	if err != nil {
		return nil, err
	}
	allResults = append(allResults, providerResults...)

	// Run credentials check
	credResults, err := r.checker.CheckCredentials()
	if err != nil {
		return nil, err
	}
	allResults = append(allResults, credResults...)

	return allResults, nil
}