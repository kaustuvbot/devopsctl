package docker

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// trivyVulnerability maps a subset of Trivy's JSON vulnerability fields.
type trivyVulnerability struct {
	VulnerabilityID string `json:"VulnerabilityID"`
	Severity        string `json:"Severity"`
	PkgName         string `json:"PkgName"`
	Title           string `json:"Title"`
}

// trivyTarget is a single Trivy result target.
type trivyTarget struct {
	Target          string               `json:"Target"`
	Vulnerabilities []trivyVulnerability `json:"Vulnerabilities"`
}

// trivyReport is the top-level structure of `trivy --format json` output.
type trivyReport struct {
	Results []trivyTarget `json:"Results"`
}

// IsTrivyInstalled returns true if the trivy binary is on PATH.
func IsTrivyInstalled() bool {
	_, err := exec.LookPath("trivy")
	return err == nil
}

// ScanImage runs trivy against an image name and returns HIGH/CRITICAL findings.
// Returns nil results (not an error) if trivy is not installed — graceful degradation.
func ScanImage(imageName string) ([]reporter.CheckResult, error) {
	if !IsTrivyInstalled() {
		return nil, nil
	}

	cmd := exec.Command("trivy", "image",
		"--format", "json",
		"--severity", "HIGH,CRITICAL",
		"--quiet",
		imageName,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("trivy scan failed for %q: %w", imageName, err)
	}

	var report trivyReport
	if err := json.Unmarshal(out, &report); err != nil {
		return nil, fmt.Errorf("failed to parse trivy output: %w", err)
	}

	var results []reporter.CheckResult
	for _, target := range report.Results {
		for _, vuln := range target.Vulnerabilities {
			if vuln.Severity != "HIGH" && vuln.Severity != "CRITICAL" {
				continue
			}
			results = append(results, reporter.CheckResult{
				CheckName:      "trivy-image-vuln",
				Severity:       vuln.Severity,
				ResourceID:     fmt.Sprintf("%s/%s", imageName, vuln.PkgName),
				Message:        fmt.Sprintf("%s: %s (%s)", vuln.VulnerabilityID, vuln.Title, vuln.PkgName),
				Recommendation: fmt.Sprintf("Update package %q to a patched version", vuln.PkgName),
			})
		}
	}
	return results, nil
}
