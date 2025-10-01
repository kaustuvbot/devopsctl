# PROJECT_SPEC.md
Project: devopsctl
Purpose: Infrastructure Hygiene & DevOps Audit Toolkit

---

# 1️⃣ Core Functional Modules

## AWS Audit

Checks:
- IAM users without MFA
- Old access keys (> configurable days)
- Users with AdministratorAccess
- Public S3 buckets
- Buckets without encryption
- Security groups open to 0.0.0.0/0 (22 or all ports)
- Unencrypted EBS volumes
- Unattached EBS volumes

Uses AWS SDK for Go v2.
Must support region override.
Must respect AWS profile.

---

## Docker Audit

Checks:
- Latest tag usage
- Missing USER directive
- Missing HEALTHCHECK
- No multi-stage build
- Risky EXPOSE usage

Optional:
- Trivy integration (--image flag)

---

## Terraform Validate

Checks:
- terraform fmt -check
- terraform validate
- Hardcoded credential detection
- Unpinned provider versions
- Optional tfsec/checkov integration

---

## Git Audit

Checks:
- Repo size threshold
- Large files (> configurable MB)
- Stale branches (> configurable days)
- Future: GitHub API integration

---

## Doctor Engine

- Aggregates results from all modules
- Weighted severity scoring
- Generates summary report
- Supports:
    --json
    --output report.md

---

# 2️⃣ Config File (.devopsctl.yaml)

Must support:

- AWS region
- Key age threshold
- Repo size threshold
- Branch age threshold
- Severity overrides

---

# 3️⃣ CLI Behavior

Commands:

devopsctl audit aws
devopsctl audit docker
devopsctl validate terraform
devopsctl audit git
devopsctl doctor

Global flags:

--json
--output <file>
--config <file>

---

# 4️⃣ Reporting Engine

Must support:

- Table output
- JSON output
- Markdown report generation
- Severity grouping
- Exit code based on highest severity

---

# 5️⃣ Extensibility

- Modules must register checks dynamically
- Adding new checks must not require CLI rewrite
- Reporter must remain centralized

---

# 6️⃣ Installation

Must support:

go install github.com/kaustuvprajapati/devopsctl/cmd/devopsctl@latest
devopsctl --help

Compatible with Ubuntu 22.04+
Go 1.21+
