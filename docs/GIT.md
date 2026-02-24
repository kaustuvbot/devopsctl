# Git Audit

The `audit git` command analyzes Git repositories for hygiene issues like oversized repos, stale branches, and large tracked files.

## Prerequisites

1. **Go 1.21+** — required for installation
2. **Git** — must be installed and available in PATH
3. **A Git repository** — the target to audit

---

## Installation

```bash
go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
```

Requires Go 1.21+. The binary is statically linked and has no runtime dependencies.

---

## Quick Start

```bash
# Audit the current repository
devopsctl audit git

# Audit a specific repository
devopsctl audit git --repo /path/to/repo

# Show only critical and high severity findings
devopsctl audit git --quiet
```

---

## Checks Reference

| Check Name | Severity | What It Flags | How to Fix |
|------------|----------|---------------|------------|
| `git-repo-size` | MEDIUM | Repository size exceeds threshold (default: 500 MB) | Use Git LFS for large files, run `git gc`, remove unnecessary objects |
| `git-stale-branch` | LOW | Branch not updated in X days (default: 90 days) | Delete stale branches: `git branch -d old-branch` |
| `git-large-file` | MEDIUM | Tracked file exceeds size threshold (default: 50 MB) | Add to `.gitignore` or use Git LFS for large files |

---

## Configuration

Create a `.devopsctl.yaml` file in your project root:

```yaml
git:
  enabled: true
  repo_size_mb: 500
  branch_age_days: 90
  large_file_mb: 50
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | true | Enable or disable Git audit |
| `repo_size_mb` | int | 500 | Flag repos larger than this size (MB) |
| `branch_age_days` | int | 90 | Flag branches older than this many days |
| `large_file_mb` | int | 50 | Flag tracked files larger than this size (MB) |

### Example: Strict Thresholds

```yaml
git:
  enabled: true
  repo_size_mb: 100      # More strict: flag repos > 100MB
  branch_age_days: 30    # More strict: flag branches > 30 days
  large_file_mb: 10      # More strict: flag files > 10MB
```

---

## Output Formats

### Table Format (default)

```
SEVERITY   CHECK NAME          RESOURCE            MESSAGE
MEDIUM     git-repo-size       /path/to/repo       Repository size exceeds threshold
LOW        git-stale-branch   old-feature         Branch has not been updated in over 90 days
MEDIUM     git-large-file     data/dump.sql       File exceeds size threshold
```

### JSON Format

```bash
devopsctl audit git --format json
```

```json
{
  "module": "git",
  "results": [
    {
      "check_name": "git-repo-size",
      "severity": "MEDIUM",
      "resource_id": "/path/to/repo",
      "message": "Repository size exceeds threshold",
      "recommendation": "Consider using Git LFS for large files or cleaning up unnecessary objects"
    }
  ]
}
```

### Markdown Format

```bash
devopsctl audit git --format markdown
```

---

## Exit Codes

| Exit Code | Meaning |
|-----------|---------|
| 0 | No findings (all checks passed) |
| 1 | LOW severity findings only |
| 2 | MEDIUM severity findings present |
| 3 | HIGH severity findings present |
| 4 | CRITICAL severity findings present |

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Git Audit
on: [push, pull_request]

jobs:
  git-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Git audit
        run: |
          go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
          devopsctl audit git
```

### GitLab CI

```yaml
git_audit:
  image: golang:1.21
  script:
    - go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
    - devopsctl audit git
  allow_failure: false
```

---

## Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "not a git repository" | Not run from a Git repo | Use `--repo` flag to specify path |
| No findings returned | Repo passes all checks | Great! Your repo follows hygiene best practices |
| Stale branch false positives | Long-running feature branches | Adjust `branch_age_days` in config |

---

## Related Commands

- `devopsctl doctor` — Run all audits including Git
- `devopsctl audit aws` — AWS infrastructure audit
- `devopsctl audit docker` — Dockerfile audit
- `devopsctl validate terraform` — Terraform validation
