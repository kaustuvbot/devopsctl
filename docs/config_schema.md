# Configuration Schema

devopsctl uses `.devopsctl.yaml` for configuration.

## File Location

The CLI looks for configuration in the following order:
1. Path specified via `--config` flag
2. `.devopsctl.yaml` in the current directory
3. `.devopsctl.yml` in the current directory
4. Default values if no config file is found

## Schema

```yaml
# AWS audit configuration
aws:
  region: us-east-1          # AWS region to scan (default: us-east-1)
  profile: default           # AWS CLI profile name (default: "")
  key_age_days: 90           # Flag access keys older than N days (default: 90)

# Docker audit configuration
docker:
  dockerfile_path: Dockerfile  # Path to Dockerfile (default: Dockerfile)

# Terraform validation configuration
terraform:
  tf_dir: .                  # Terraform directory (default: "")

# Git audit configuration
git:
  repo_size_mb: 500          # Flag repos larger than N MB (default: 500)
  branch_age_days: 90        # Flag branches older than N days (default: 90)
  large_file_mb: 50          # Flag files larger than N MB (default: 50)
```

## Default Values

When a configuration file is not found or a field is omitted,
the following defaults are used:

| Field | Default |
|-------|---------|
| `aws.region` | `us-east-1` |
| `aws.key_age_days` | `90` |
| `docker.dockerfile_path` | `Dockerfile` |
| `git.repo_size_mb` | `500` |
| `git.branch_age_days` | `90` |
| `git.large_file_mb` | `50` |

## Validation

Invalid YAML will cause devopsctl to exit with an error.
Missing optional fields will use default values.
