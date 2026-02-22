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
| Terraform Validation | [docs/TERRAFORM.md](docs/TERRAFORM.md) |

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

## Global Flags

```
--format <fmt>   Output format: table, json, markdown (default: table)
--quiet          Show only CRITICAL and HIGH severity findings
--output <file>  Write report to file
--config <file>  Path to config file (default: .devopsctl.yaml)
--json           Output in JSON format (deprecated, use --format json)
```

## Requirements

- Go 1.21+
- Ubuntu Linux (primary target)

## License

MIT
