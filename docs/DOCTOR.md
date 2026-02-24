# Doctor

The `doctor` command runs all audit and validation modules and produces an aggregated health report. It's the flagship command for comprehensive infrastructure hygiene checks.

## Overview

Doctor orchestrates all registered modules:
- AWS audit (IAM, S3, EC2, EBS)
- Docker audit (Dockerfile checks)
- Terraform validation (fmt, validate, credentials)
- Git audit (repo size, branches, large files)

Each module runs independently, and results are aggregated into a summary with severity scoring.

---

## Installation

```bash
go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
```

Requires Go 1.21+. The binary is statically linked and has no runtime dependencies.

---

## Quick Start

```bash
# Run all audits and generate health report
devopsctl doctor

# Show only critical and high severity findings
devopsctl doctor --quiet

# Output as JSON for automation
devopsctl doctor --format json

# Write report to file
devopsctl doctor --output report.md --format markdown
```

---

## Severity Scoring

Doctor computes a weighted score based on findings:

| Severity | Weight |
|----------|--------|
| LOW | 1 |
| MEDIUM | 2 |
| HIGH | 3 |
| CRITICAL | 4 |

The total score is the sum of all finding weights. Use the score to track hygiene trends over time.

---

## Output Summary

Doctor output includes:
- **Module Results** — individual findings from each module
- **Summary** — total findings, breakdown by severity, weighted score
- **Module Errors** — any modules that failed to run

### Table Format Example

```
=== AWS Audit ===
SEVERITY   CHECK NAME              RESOURCE      MESSAGE
CRITICAL   s3-public-bucket       my-bucket     Bucket is publicly accessible

=== Docker Audit ===
No issues found.

=== Terraform Validation ===
No issues found.

=== Git Audit ===
No issues found.

=== Summary ===
Total: 1 | Critical: 1 | High: 0 | Medium: 0 | Low: 0
Score: 4
```

---

## Output Formats

### Table Format (default)

Human-readable format with color coding:
- CRITICAL: red
- HIGH: yellow
- MEDIUM: green
- LOW: green

### JSON Format

```bash
devopsctl doctor --format json
```

```json
{
  "doctor": true,
  "summary": {
    "total_findings": 1,
    "critical": 1,
    "high": 0,
    "medium": 0,
    "low": 0,
    "score": 4,
    "modules_failed": 0,
    "module_errors": {}
  },
  "reports": [
    {
      "module": "aws",
      "results": [
        {
          "check_name": "s3-public-bucket",
          "severity": "CRITICAL",
          "resource_id": "my-bucket",
          "message": "Bucket is publicly accessible",
          "recommendation": "Block public access in S3 settings"
        }
      ]
    }
  ]
}
```

### Markdown Format

```bash
devopsctl doctor --format markdown
```

Generates a Markdown report suitable for:
- Team sharing
- CI/CD artifacts
- Documentation

---

## Exit Codes

Doctor exits with the highest severity found across all modules:

| Exit Code | Meaning |
|-----------|---------|
| 0 | No findings (all modules passed) |
| 1 | LOW severity findings only |
| 2 | MEDIUM severity findings present |
| 3 | HIGH severity findings present |
| 4 | CRITICAL severity findings present |

---

## Module Toggle

Disable specific modules in `.devopsctl.yaml`:

```yaml
aws:
  enabled: true
docker:
  enabled: false     # Skip Docker audit
terraform:
  enabled: true
git:
  enabled: true
```

When a module is disabled, it's skipped in `doctor` runs.

---

## Ignoring Specific Checks

Ignore specific checks across all modules:

```yaml
ignore:
  checks:
    - dockerfile-no-healthcheck
    - git-stale-branch
```

Checks are matched exactly by name.

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Infrastructure Health Check
on:
  schedule:
    - cron: '0 0 * * *'  # Daily
  push:
    branches: [main]
  pull_request:

jobs:
  doctor:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run doctor
        run: |
          go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
          devopsctl doctor --format json --output doctor-report.json

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: health-report
          path: doctor-report.json
```

### GitLab CI

```yaml
doctor:
  image: golang:1.21
  script:
    - go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
    - devopsctl doctor --format markdown --output health-report.md
  artifacts:
    paths:
      - health-report.md
```

---

## Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Module fails to run | Missing credentials or dependencies | Check prerequisites for each module |
| Partial results | One module failed but others succeeded | Check `module_errors` in JSON output |
| No findings | All modules passed | Great! Your infrastructure is healthy |

---

## Related Commands

- `devopsctl audit aws` — Run only AWS audit
- `devopsctl audit docker` — Run only Docker audit
- `devopsctl audit git` — Run only Git audit
- `devopsctl validate terraform` — Run only Terraform validation
