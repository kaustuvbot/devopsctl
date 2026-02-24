# devopsctl

A lightweight CLI toolkit for infrastructure hygiene and DevOps validation.

## Features

- AWS infrastructure audit (IAM, S3, EC2, EBS)
- Docker security analysis (Dockerfile checks, optional Trivy integration)
- Terraform validation (fmt, validate, credential detection)
- Git repository hygiene checks (repo size, large files, stale branches)
- Aggregated health reports via doctor engine

## Documentation

| Module | Guide |
|--------|-------|
| AWS Audit | [docs/AWS.md](docs/AWS.md) |
| Docker Audit | [docs/DOCKER.md](docs/DOCKER.md) |
| Terraform Validation | [docs/TERRAFORM.md](docs/TERRAFORM.md) |
| Git Audit | [docs/GIT.md](docs/GIT.md) |
| Doctor | [docs/DOCTOR.md](docs/DOCTOR.md) |
| Configuration | [docs/CONFIGURATION.md](docs/CONFIGURATION.md) |

## Installation

```bash
go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
```

## Usage

```bash
devopsctl --help
devopsctl audit aws
devopsctl audit docker --file Dockerfile
devopsctl validate terraform
devopsctl audit git
devopsctl doctor
```

## Configuration

Create a `.devopsctl.yaml` file in your project root:

```yaml
aws:
  enabled: true
  region: us-east-1
  key_age_days: 90

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

ignore:
  checks: []
```

See [docs/CONFIGURATION.md](docs/CONFIGURATION.md) for full configuration options.

## Global Flags

```
--format <fmt>   Output format: table, json, markdown (default: table)
--quiet          Show only CRITICAL and HIGH severity findings
--output <file>  Write report to file
--config <file>  Path to config file (default: .devopsctl.yaml)
--json           Output in JSON format (deprecated, use --format json)
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | No findings |
| 1 | LOW severity |
| 2 | MEDIUM severity |
| 3 | HIGH severity |
| 4 | CRITICAL severity |

## Requirements

- Go 1.21+
- Ubuntu Linux (primary target)

## License

MIT
