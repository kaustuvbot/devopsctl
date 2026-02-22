package git

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	appconfig "github.com/kaustuvbot/devopsctl/internal/config"
	"github.com/kaustuvbot/devopsctl/internal/reporter"
	"github.com/kaustuvbot/devopsctl/internal/severity"
)

// CheckLargeFiles checks for tracked files exceeding the configured size threshold.
// Returns MEDIUM severity findings for each large file.
func CheckLargeFiles(ctx context.Context, client *Client, cfg appconfig.GitConfig) ([]reporter.CheckResult, error) {
	output, err := client.Run(ctx, "ls-files", "-s")
	if err != nil {
		return nil, err
	}

	var results []reporter.CheckResult
	thresholdKB := cfg.LargeFileMB * 1024

	// Parse: <mode> <object> <stage> <filename>
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		size, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		// Size is in bytes, convert to KB for comparison
		sizeKB := size / 1024

		if sizeKB > thresholdKB {
			filename := strings.Join(parts[3:], " ")
			results = append(results, reporter.CheckResult{
				CheckName:      "git-large-file",
				Severity:       string(severity.Medium),
				ResourceID:     filename,
				Message:        "File exceeds size threshold",
				Recommendation: "Consider using Git LFS for large files or removing from version control",
			})
		}
	}

	return results, nil
}

// CheckLargeFilesRegex uses regex pattern matching to find large files.
// Alternative implementation for more flexible matching.
func CheckLargeFilesRegex(ctx context.Context, client *Client, cfg appconfig.GitConfig, pattern string) ([]reporter.CheckResult, error) {
	output, err := client.Run(ctx, "ls-files", "-z")
	if err != nil {
		return nil, err
	}

	var results []reporter.CheckResult
	re := regexp.MustCompile(pattern)

	// Files are null-terminated
	files := strings.Split(output, "\x00")
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" || !re.MatchString(file) {
			continue
		}

		// Get size for matching file
		sizeOutput, err := client.Run(ctx, "ls-files", "-s", file)
		if err != nil {
			continue
		}

		parts := strings.Fields(sizeOutput)
		if len(parts) < 2 {
			continue
		}

		size, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		sizeKB := size / 1024
		if sizeKB > cfg.LargeFileMB*1024 {
			results = append(results, reporter.CheckResult{
				CheckName:      "git-large-file",
				Severity:       string(severity.Medium),
				ResourceID:     file,
				Message:        "File matches pattern and exceeds size threshold",
				Recommendation: "Consider using Git LFS or removing from version control",
			})
		}
	}

	return results, nil
}
