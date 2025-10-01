# devopsctl

A lightweight CLI toolkit for infrastructure hygiene and DevOps validation.

## Features

- AWS infrastructure audit (IAM, S3, EC2, EBS)
- Docker security analysis (Dockerfile checks, optional Trivy integration)
- Terraform validation (fmt, validate, credential detection)
- Git repository hygiene checks (repo size, large files, stale branches)
- Aggregated health reports via doctor engine

## Installation

```bash
go install github.com/kaustuvprajapati/devopsctl/cmd/devopsctl@latest
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
--json           Output in JSON format
--output <file>  Write report to file
--config <file>  Path to config file (default: .devopsctl.yaml)
```

## Requirements

- Go 1.21+
- Ubuntu Linux (primary target)

## License

MIT
