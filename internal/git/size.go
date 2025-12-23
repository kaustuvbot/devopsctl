package git

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	appconfig "github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
	"github.com/kaustuvprajapati/devopsctl/internal/severity"
)

// CheckRepoSize checks if the repository size exceeds the configured threshold.
// Returns a MEDIUM severity finding if the repo is too large.
func CheckRepoSize(ctx context.Context, client *Client, cfg appconfig.GitConfig) ([]reporter.CheckResult, error) {
	output, err := client.Run(ctx, "count-objects", "-vH")
	if err != nil {
		return nil, err
	}

	// Parse output like:
	// count: 1234
	// size: 123.45MiB
	// in-pack: 5678
	re := regexp.MustCompile(`size:\s+(\d+\.?\d*)([KMGT]i?B)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) < 3 {
		return nil, nil // Cannot determine size, skip
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return nil, err
	}

	unit := matches[2]
	sizeMB := convertToMB(value, unit)

	if sizeMB > float64(cfg.RepoSizeMB) {
		return []reporter.CheckResult{
			{
				CheckName:      "git-repo-size",
				Severity:       string(severity.Medium),
				ResourceID:     client.repoPath,
				Message:        "Repository size exceeds threshold",
				Recommendation: "Consider using Git LFS for large files or cleaning up unnecessary objects",
			},
		}, nil
	}

	return nil, nil
}

// convertToMB converts a size value to megabytes.
func convertToMB(value float64, unit string) float64 {
	unitUpper := strings.ToUpper(unit)
	switch unitUpper {
	case "B":
		return value / (1024 * 1024)
	case "KB", "KIB":
		return value / 1024
	case "MB", "MIB":
		return value
	case "GB", "GIB":
		return value * 1024
	case "TB", "TIB":
		return value * 1024 * 1024
	default:
		return value
	}
}
