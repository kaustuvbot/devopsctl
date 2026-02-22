package docker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// CheckLatestTag detects FROM instructions using :latest or untagged images.
// Using mutable tags makes builds non-reproducible.
// Severity: MEDIUM
func CheckLatestTag(df *ParsedDockerfile) []reporter.CheckResult {
	var results []reporter.CheckResult
	for _, instr := range df.Instructions {
		if instr.Command != "FROM" {
			continue
		}
		// Strip AS alias: "FROM ubuntu:latest AS builder" -> "ubuntu:latest"
		image := strings.Fields(instr.Args)[0]
		// Skip scratch
		if image == "scratch" {
			continue
		}
		// Untagged (no colon, no digest) defaults to :latest
		hasTag := strings.Contains(image, ":")
		hasDigest := strings.Contains(image, "@")
		if !hasTag && !hasDigest {
			results = append(results, reporter.CheckResult{
				CheckName:      "dockerfile-latest-tag",
				Severity:       "MEDIUM",
				ResourceID:     fmt.Sprintf("%s:line%d", df.Path, instr.LineNum),
				Message:        fmt.Sprintf("FROM uses untagged image %q (defaults to :latest) at line %d", image, instr.LineNum),
				Recommendation: "Pin image to a specific digest or immutable tag (e.g., ubuntu:22.04)",
			})
			continue
		}
		if hasTag {
			parts := strings.SplitN(image, ":", 2)
			tag := parts[1]
			if tag == "latest" {
				results = append(results, reporter.CheckResult{
					CheckName:      "dockerfile-latest-tag",
					Severity:       "MEDIUM",
					ResourceID:     fmt.Sprintf("%s:line%d", df.Path, instr.LineNum),
					Message:        fmt.Sprintf("FROM uses mutable :latest tag: %q at line %d", image, instr.LineNum),
					Recommendation: "Pin image to a specific digest or immutable tag (e.g., ubuntu:22.04)",
				})
			}
		}
	}
	return results
}

// CheckNoUser detects Dockerfiles that never set a non-root USER directive.
// Containers running as root increase the blast radius of a container escape.
// Severity: HIGH
func CheckNoUser(df *ParsedDockerfile) []reporter.CheckResult {
	for _, instr := range df.Instructions {
		if instr.Command == "USER" {
			user := strings.TrimSpace(instr.Args)
			if user != "0" && user != "root" {
				return nil // non-root USER found
			}
		}
	}
	return []reporter.CheckResult{{
		CheckName:      "dockerfile-runs-as-root",
		Severity:       "HIGH",
		ResourceID:     df.Path,
		Message:        "Dockerfile has no USER directive; container will run as root",
		Recommendation: "Add a USER directive with a non-root user (e.g., USER 1001)",
	}}
}

// CheckNoHealthcheck detects Dockerfiles missing a HEALTHCHECK instruction.
// Without a healthcheck, orchestrators cannot detect unhealthy containers.
// Severity: LOW
func CheckNoHealthcheck(df *ParsedDockerfile) []reporter.CheckResult {
	for _, instr := range df.Instructions {
		if instr.Command == "HEALTHCHECK" {
			return nil
		}
	}
	return []reporter.CheckResult{{
		CheckName:      "dockerfile-no-healthcheck",
		Severity:       "LOW",
		ResourceID:     df.Path,
		Message:        "Dockerfile has no HEALTHCHECK instruction",
		Recommendation: "Add HEALTHCHECK to allow container orchestrators to monitor service health",
	}}
}

// CheckNoMultiStage detects Dockerfiles with only a single FROM stage.
// Multi-stage builds reduce final image size and remove build tools from production.
// Severity: LOW
func CheckNoMultiStage(df *ParsedDockerfile) []reporter.CheckResult {
	fromCount := 0
	for _, instr := range df.Instructions {
		if instr.Command == "FROM" {
			fromCount++
		}
	}
	if fromCount < 2 {
		return []reporter.CheckResult{{
			CheckName:      "dockerfile-no-multi-stage",
			Severity:       "LOW",
			ResourceID:     df.Path,
			Message:        "Dockerfile uses a single-stage build",
			Recommendation: "Consider multi-stage builds to reduce final image size and exclude build tools",
		}}
	}
	return nil
}

// riskyPorts maps sensitive port numbers to their service names.
var riskyPorts = map[int]string{
	22:    "SSH",
	23:    "Telnet",
	3306:  "MySQL",
	5432:  "PostgreSQL",
	6379:  "Redis",
	27017: "MongoDB",
}

// CheckRiskyExpose detects EXPOSE of sensitive or privileged ports.
// Severity: MEDIUM
func CheckRiskyExpose(df *ParsedDockerfile) []reporter.CheckResult {
	var results []reporter.CheckResult
	for _, instr := range df.Instructions {
		if instr.Command != "EXPOSE" {
			continue
		}
		for _, portStr := range strings.Fields(instr.Args) {
			// Strip protocol suffix: "22/tcp" -> "22"
			portStr = strings.Split(portStr, "/")[0]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}
			if name, risky := riskyPorts[port]; risky {
				results = append(results, reporter.CheckResult{
					CheckName:      "dockerfile-risky-expose",
					Severity:       "MEDIUM",
					ResourceID:     fmt.Sprintf("%s:line%d", df.Path, instr.LineNum),
					Message:        fmt.Sprintf("EXPOSE includes risky port %d (%s) at line %d", port, name, instr.LineNum),
					Recommendation: fmt.Sprintf("Avoid exposing sensitive service port %d unless intentional", port),
				})
			}
		}
	}
	return results
}
