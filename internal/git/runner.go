package git

import (
	"context"
	"fmt"

	appconfig "github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// Runner orchestrates git audit checks.
type Runner struct {
	client *Client
	cfg    appconfig.GitConfig
}

// NewRunner creates a new git runner.
func NewRunner(repoPath string, cfg appconfig.GitConfig) *Runner {
	return &Runner{
		client: NewClient(repoPath),
		cfg:    cfg,
	}
}

// RunAll executes all git checks and returns aggregated results.
func (r *Runner) RunAll(ctx context.Context) ([]reporter.CheckResult, error) {
	var all []reporter.CheckResult
	var errs []string

	checks := []struct {
		name string
		fn   func(context.Context) ([]reporter.CheckResult, error)
	}{
		{"repo-size", func(ctx context.Context) ([]reporter.CheckResult, error) {
			return CheckRepoSize(ctx, r.client, r.cfg)
		}},
		{"stale-branches", func(ctx context.Context) ([]reporter.CheckResult, error) {
			return CheckStaleBranches(ctx, r.client, r.cfg)
		}},
		{"large-files", func(ctx context.Context) ([]reporter.CheckResult, error) {
			return CheckLargeFiles(ctx, r.client, r.cfg)
		}},
	}

	for _, check := range checks {
		results, err := check.fn(ctx)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", check.name, err))
			continue
		}
		all = append(all, results...)
	}

	if len(errs) > 0 {
		return all, fmt.Errorf("some checks failed: %v", errs)
	}

	return all, nil
}

// RunAllSimple is a convenience method without context.
func (r *Runner) RunAllSimple() ([]reporter.CheckResult, error) {
	return r.RunAll(context.Background())
}
