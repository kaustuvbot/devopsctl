package docker

import (
	"fmt"

	appconfig "github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// RunOptions controls optional Docker audit behaviors.
type RunOptions struct {
	ImageName string // If set, also run Trivy scan against this image
}

// RunAll executes all Dockerfile static checks and optional Trivy scan.
// Returns all findings aggregated into a single []reporter.CheckResult.
// Checks that fail are skipped, not fatal - returns partial results.
func RunAll(cfg appconfig.DockerConfig, opts RunOptions) ([]reporter.CheckResult, error) {
	var all []reporter.CheckResult
	var errs []string

	// Parse Dockerfile first - this is a prerequisite check
	df, err := ParseDockerfile(cfg.DockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("docker audit: %w", err)
	}

	// Run static checks - collect errors, continue on failure
	type checkFn func() []reporter.CheckResult
	checks := []checkFn{
		func() []reporter.CheckResult { return CheckLatestTag(df) },
		func() []reporter.CheckResult { return CheckNoUser(df) },
		func() []reporter.CheckResult { return CheckNoHealthcheck(df) },
		func() []reporter.CheckResult { return CheckNoMultiStage(df) },
		func() []reporter.CheckResult { return CheckRiskyExpose(df) },
	}

	for _, check := range checks {
		results := check()
		all = append(all, results...)
	}

	// Optional Trivy scan
	if opts.ImageName != "" {
		trivyResults, err := ScanImage(opts.ImageName)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			all = append(all, trivyResults...)
		}
	}

	if len(errs) > 0 {
		return all, fmt.Errorf("some checks failed: %v", errs)
	}
	return all, nil
}
