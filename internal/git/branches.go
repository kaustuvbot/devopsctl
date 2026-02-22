package git

import (
	"context"
	"strconv"
	"strings"
	"time"

	appconfig "github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

// CheckStaleBranches checks if any branches have not been updated in the configured number of days.
// Returns LOW severity findings for stale branches.
func CheckStaleBranches(ctx context.Context, client *Client, cfg appconfig.GitConfig) ([]reporter.CheckResult, error) {
	output, err := client.Run(ctx, "branch", "-a", "--format=%(refname:short)|%(committerdate:unix)")
	if err != nil {
		return nil, err
	}

	var results []reporter.CheckResult
	threshold := time.Duration(cfg.BranchAgeDays) * 24 * time.Hour
	now := time.Now()

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip HEAD references
		if strings.Contains(line, "HEAD") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		branch := strings.TrimSpace(parts[0])
		// Skip remote tracking branch prefixes
		branch = strings.TrimPrefix(branch, "remotes/origin/")

		timestamp := strings.TrimSpace(parts[1])
		unixTime, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			continue
		}

		commitTime := time.Unix(unixTime, 0)
		age := now.Sub(commitTime)

		if age > threshold {
			results = append(results, reporter.CheckResult{
				CheckName:      "git-stale-branch",
				Severity:       string(severity.Low),
				ResourceID:     branch,
				Message:        "Branch has not been updated in over " + string(rune(cfg.BranchAgeDays)) + " days",
				Recommendation: "Consider deleting stale branches or merging/updating them",
			})
		}
	}

	return results, nil
}
