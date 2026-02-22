package aws

import (
	"context"
	"fmt"

	appconfig "github.com/kaustuvbot/devopsctl/internal/config"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// RunAll executes all AWS checks and returns aggregated results.
// Checks that fail due to insufficient permissions are skipped, not fatal.
func RunAll(ctx context.Context, clients *AWSClients, cfg appconfig.AWSConfig) ([]reporter.CheckResult, error) {
	var all []reporter.CheckResult
	var errs []string

	type checkFn func() ([]reporter.CheckResult, error)
	checks := []checkFn{
		func() ([]reporter.CheckResult, error) { return CheckIAMUsersMFA(ctx, clients.IAM) },
		func() ([]reporter.CheckResult, error) { return CheckIAMAccessKeyAge(ctx, clients.IAM, cfg.KeyAgeDays) },
		func() ([]reporter.CheckResult, error) { return CheckIAMAdminUsers(ctx, clients.IAM) },
		func() ([]reporter.CheckResult, error) { return CheckS3PublicBuckets(ctx, clients.S3) },
		func() ([]reporter.CheckResult, error) { return CheckS3Encryption(ctx, clients.S3) },
		func() ([]reporter.CheckResult, error) { return CheckS3Versioning(ctx, clients.S3) },
		func() ([]reporter.CheckResult, error) { return CheckSecurityGroups(ctx, clients.EC2) },
		func() ([]reporter.CheckResult, error) { return CheckEBSEncryption(ctx, clients.EC2) },
		func() ([]reporter.CheckResult, error) { return CheckEBSUnattached(ctx, clients.EC2) },
	}

	for _, check := range checks {
		results, err := check()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		all = append(all, results...)
	}

	if len(errs) > 0 {
		return all, fmt.Errorf("some checks failed: %v", errs)
	}
	return all, nil
}
