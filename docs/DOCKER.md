# Docker Audit

The `audit docker` command analyzes Dockerfiles for security risks, best practices, and common misconfigurations. It performs static analysis without requiring a Docker build.

## Prerequisites

1. **Go 1.21+** — required for installation
2. **A Dockerfile** — the target file to audit

No Docker daemon is required — the tool performs static analysis on the Dockerfile text.

---

## Installation

```bash
go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
```

Requires Go 1.21+. The binary is statically linked and has no runtime dependencies.

---

## Quick Start

```bash
# Audit the Dockerfile in the current directory
devopsctl audit docker

# Audit a specific Dockerfile
devopsctl audit docker --file /path/to/Dockerfile

# Show only critical and high severity findings
devopsctl audit docker --quiet
```

---

## Checks Reference

| Check Name | Severity | What It Flags | How to Fix |
|------------|----------|---------------|------------|
| `dockerfile-latest-tag` | MEDIUM | Base image uses `:latest` tag or no tag (defaults to :latest) | Pin to specific version: `FROM ubuntu:22.04` |
| `dockerfile-runs-as-root` | HIGH | No USER directive — container runs as root | Add non-root user: `USER 1001` or `USER app` |
| `dockerfile-no-healthcheck` | LOW | No HEALTHCHECK instruction | Add health check: `HEALTHCHECK CMD curl -f http://localhost/ || exit 1` |
| `dockerfile-no-multi-stage` | LOW | Single-stage build | Use multi-stage: `FROM builder AS builder` then `COPY --from=builder` |
| `dockerfile-risky-expose` | MEDIUM | Exposes sensitive port (22, 23, 3306, 5432, 6379, 27017) | Remove or restrict EXPOSE unless intentional |

---

## Configuration

Create a `.devopsctl.yaml` file in your project root:

```yaml
docker:
  enabled: true
  dockerfile_path: Dockerfile
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | true | Enable or disable Docker audit |
| `dockerfile_path` | string | "Dockerfile" | Path to the Dockerfile to audit |

---

## Output Formats

### Table Format (default)

```
SEVERITY   CHECK NAME                    RESOURCE      MESSAGE
MEDIUM     dockerfile-latest-tag          Dockerfile:5  FROM uses mutable :latest tag
HIGH       dockerfile-runs-as-root        Dockerfile    No USER directive; container runs as root
```

### JSON Format

```bash
devopsctl audit docker --format json
```

```json
{
  "module": "docker",
  "results": [
    {
      "check_name": "dockerfile-latest-tag",
      "severity": "MEDIUM",
      "resource_id": "Dockerfile:5",
      "message": "FROM uses mutable :latest tag: ubuntu:latest at line 5",
      "recommendation": "Pin image to a specific digest or immutable tag (e.g., ubuntu:22.04)"
    }
  ]
}
```

### Markdown Format

```bash
devopsctl audit docker --format markdown
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

## Trivy Integration (Optional)

For container vulnerability scanning, install [Trivy](https://aquasecurity.github.io/trivy/) separately:

```bash
# Install Trivy
curl -sfL https://aquasecurity.github.io/trivy/install.sh | sh

# Scan a built image (requires Docker)
devopsctl audit docker --image myapp:latest
```

> **Note**: Trivy integration requires Docker to be installed and the image to be built. This is optional — the core `audit docker` command works without Trivy.

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Docker Audit
on: [push, pull_request]

jobs:
  docker-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Docker audit
        run: |
          go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
          devopsctl audit docker --file Dockerfile
```

### GitLab CI

```yaml
docker_audit:
  image: golang:1.21
  script:
    - go install github.com/kaustuvbot/devopsctl/cmd/devopsctl@latest
    - devopsctl audit docker --file Dockerfile
  allow_failure: false
```

---

## Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Dockerfile not found" | File not at expected path | Use `--file` flag or update `dockerfile_path` in config |
| No findings returned | Dockerfile passes all checks | Great! Your Dockerfile follows best practices |
| Parser error | Unrecognized Dockerfile syntax | Ensure valid Dockerfile syntax |

---

## Related Commands

- `devopsctl doctor` — Run all audits including Docker
- `devopsctl audit aws` — AWS infrastructure audit
- `devopsctl validate terraform` — Terraform validation
- `devopsctl audit git` — Git repository hygiene
