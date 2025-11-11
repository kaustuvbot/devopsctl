package docker

import (
	"fmt"
	"os"

	appconfig "github.com/kaustuvprajapati/devopsctl/internal/config"
	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// RunOptions controls optional Docker audit behaviors.
type RunOptions struct {
	ImageName string // If set, also run Trivy scan against this image
}

// RunAll executes all Dockerfile static checks and optional Trivy scan.
// Returns all findings aggregated into a single []reporter.CheckResult.
func RunAll(cfg appconfig.DockerConfig, opts RunOptions) ([]reporter.CheckResult, error) {
	df, err := ParseDockerfile(cfg.DockerfilePath)
	if err != nil {
		return nil, fmt.Errorf("docker audit: %w", err)
	}

	var all []reporter.CheckResult
	all = append(all, CheckLatestTag(df)...)
	all = append(all, CheckNoUser(df)...)
	all = append(all, CheckNoHealthcheck(df)...)
	all = append(all, CheckNoMultiStage(df)...)
	all = append(all, CheckRiskyExpose(df)...)

	if opts.ImageName != "" {
		trivyResults, err := ScanImage(opts.ImageName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: trivy scan failed: %v\n", err)
		} else {
			all = append(all, trivyResults...)
		}
	}

	return all, nil
}
