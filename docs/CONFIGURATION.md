# Configuration

devopsctl reads configuration from a YAML file to customize behavior across all modules.

## File Location

The tool searches for configuration in this order:

1. **Flag**: Path specified via `--config` flag
   ```bash
   devopsctl --config /path/to/custom.yaml doctor
   ```

2. **Current directory**: `.devopsctl.yaml` or `.devopsctl.yml`

3. **Default values**: If no config file is found, sensible defaults are used

---

## Full Schema

```yaml
# AWS Audit Configuration
aws:
  enabled: true           # Enable/disable module
  region: us-east-1       # AWS region to scan
  profile: ""             # AWS CLI profile (empty = default chain)
  key_age_days: 90        # Flag access keys older than N days

# Docker Audit Configuration
docker:
  enabled: true
  dockerfile_path: Dockerfile

# Terraform Validation Configuration
terraform:
  enabled: true
  tf_dir: .

# Git Audit Configuration
git:
  enabled: true
  repo_size_mb: 500       # Flag repos larger than N MB
  branch_age_days: 90     # Flag branches older than N days
  large_file_mb: 50       # Flag files larger than N MB

# Ignore specific checks
ignore:
  checks:
    - check-name-1       # Exact match on check name
    - check-name-2
```

---

## Module Toggle

Disable specific modules by setting `enabled: false`:

```yaml
aws:
  enabled: true
docker:
  enabled: false    # Skip Docker audit in doctor command
terraform:
  enabled: true
git:
  enabled: true
```

When a module is disabled:
- Running `devopsctl audit <module>` still works
- Running `devopsctl doctor` skips the disabled module

---

## Ignore Specific Checks

The `ignore.checks` list filters out specific findings across all modules:

```yaml
ignore:
  checks:
    - dockerfile-no-healthcheck
    - git-stale-branch
```

Checks are matched by exact name. The finding is removed from results before severity filtering.

---

## AWS Configuration

### Region

```yaml
aws:
  region: us-west-2
```

Default: `us-east-1`

### Profile

```yaml
aws:
  profile: my-profile
```

When empty, uses the default credential chain (env vars → config file → EC2 role).

### Key Age Threshold

```yaml
aws:
  key_age_days: 60
```

Flags IAM access keys older than N days. Default: 90 days.

---

## Docker Configuration

### Dockerfile Path

```yaml
docker:
  dockerfile_path: containers/app/Dockerfile
```

Default: `Dockerfile` in current directory

---

## Terraform Configuration

### Working Directory

```yaml
terraform:
  tf_dir: ./infrastructure
```

Default: `.` (current directory)

---

## Git Configuration

### Repository Size Threshold

```yaml
git:
  repo_size_mb: 1000
```

Default: 500 MB

### Branch Age Threshold

```yaml
git:
  branch_age_days: 180
```

Default: 90 days

### Large File Threshold

```yaml
git:
  large_file_mb: 100
```

Default: 50 MB

---

## Examples

### Minimal Config (Defaults)

```yaml
# Uses all defaults - equivalent to no config file
```

### Strict Thresholds

```yaml
aws:
  enabled: true
  key_age_days: 30    # Stricter: flag keys > 30 days

docker:
  enabled: true

terraform:
  enabled: true

git:
  enabled: true
  repo_size_mb: 100   # Stricter: flag repos > 100MB
  branch_age_days: 30 # Stricter: flag branches > 30 days
  large_file_mb: 10   # Stricter: flag files > 10MB

ignore:
  checks: []
```

### CI/CD Optimized

```yaml
aws:
  enabled: true
  region: us-east-1

docker:
  enabled: true
  dockerfile_path: Dockerfile

terraform:
  enabled: true
  tf_dir: .

git:
  enabled: true
  repo_size_mb: 500
  branch_age_days: 90
  large_file_mb: 50

# Ignore known acceptable findings
ignore:
  checks:
    - dockerfile-no-healthcheck    # Optional in some contexts
```

### Disable All Except AWS

```yaml
aws:
  enabled: true
docker:
  enabled: false
terraform:
  enabled: false
git:
  enabled: false
```

---

## Config File Search Order

1. `--config` flag (highest priority)
2. `.devopsctl.yaml` in current directory
3. `.devopsctl.yml` in current directory
4. Use defaults (if no file found)

---

## Validation

Config is validated on load. Common errors:

| Error | Cause | Fix |
|-------|-------|-----|
| Unknown field | Typo in field name | Check YAML indentation and spelling |
| Invalid type | Wrong value type | Ensure numbers are numbers, booleans are booleans |
| File not found | Path doesn't exist | Check `--config` path or current directory |

---

## Related Commands

- `devopsctl doctor` — Run all modules with config
- `devopsctl audit aws` — AWS audit with config
- `devopsctl audit docker` — Docker audit with config
- `devopsctl audit git` — Git audit with config
- `devopsctl validate terraform` — Terraform validation with config
