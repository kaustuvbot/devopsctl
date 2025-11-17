# devopsctl Usage Guide

## Installation

```bash
go install github.com/kaustuvprajapati/devopsctl/cmd/devopsctl@latest
```

Or build locally:

```bash
git clone https://github.com/kaustuvprajapati/devopsctl
cd devopsctl
go build -o devopsctl ./cmd/devopsctl
```

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--config <path>` | Config file path (default: `.devopsctl.yaml`) |
| `--json` | Output results in JSON format |
| `--output <file>` | Write report to a file |

---

## Commands

### `audit` — Infrastructure Audit Checks

```bash
devopsctl audit aws        # Audit AWS IAM, S3, EC2
devopsctl audit docker     # Audit Dockerfile
devopsctl audit git        # Audit Git repository hygiene
```

### `validate` — Validate Infrastructure Code

```bash
devopsctl validate terraform   # Validate Terraform configuration
```

### `doctor` — Full Health Report

Runs all checks and aggregates results into a single report.

```bash
devopsctl doctor
devopsctl doctor --json
devopsctl doctor --output report.md
```

### `version` — Print Version

```bash
devopsctl version
```

---

## Configuration

Create `.devopsctl.yaml` in your project root:

```yaml
aws:
  region: us-east-1
  profile: default
  key_age_days: 90

docker:
  dockerfile_path: Dockerfile

terraform:
  tf_dir: .

git:
  repo_size_mb: 500
  branch_age_days: 90
  large_file_mb: 50
```

All fields are optional — omitted fields use the defaults shown above.
See [config_schema.md](config_schema.md) for full reference.

---

## Output Format

Each check produces a `CheckResult`:

| Field | Description |
|-------|-------------|
| `check_name` | Name of the check |
| `severity` | `LOW`, `MEDIUM`, `HIGH`, or `CRITICAL` |
| `resource_id` | Affected resource identifier |
| `message` | What was found |
| `recommendation` | How to fix it |

Exit codes reflect the highest severity found: `0` (clean), `1` (LOW), `2` (MEDIUM), `3` (HIGH), `4` (CRITICAL).

---

## Internal Packages

| Package | Purpose |
|---------|---------|
| `internal/cli` | Cobra command definitions |
| `internal/config` | YAML config loading and defaults |
| `internal/reporter` | Output formatting (JSON, table) |
| `internal/severity` | Severity levels, weights, exit codes |

Module-specific packages (`aws`, `docker`, `terraform`, `git`, `doctor`) are added in subsequent batches.

---

## Further Reading

- [config_schema.md](config_schema.md) — Full config reference
- [aws_audit.md](aws_audit.md) — AWS module spec
- [docker_audit.md](docker_audit.md) — Docker module spec
- [terraform_validate.md](terraform_validate.md) — Terraform module spec
- [git_audit.md](git_audit.md) — Git module spec
- [doctor_engine.md](doctor_engine.md) — Doctor engine spec
