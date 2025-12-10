# Terraform Validation

The `validate terraform` command runs automated checks on your Terraform configuration to detect formatting issues, configuration errors, and security concerns.

## Quick Start

```bash
# Validate Terraform in current directory
devopsctl validate terraform

# Validate specific directory
devopsctl validate terraform --dir ./infrastructure
```

## Checks Performed

### 🔧 Format Check
Ensures all Terraform files follow canonical formatting standards.

```bash
# What it does
terraform fmt -check -recursive

# If it fails
devopsctl validate terraform
# Shows: "Terraform files are not properly formatted"
# Fix it: terraform fmt -r .
```

### ✔️ Validation Check
Verifies that Terraform configuration is syntactically valid and has no configuration errors.

```bash
# What it does
terraform validate

# If it fails
devopsctl validate terraform
# Shows: "Terraform configuration is invalid"
# Check: terraform validate (for detailed error messages)
```

### 📌 Provider Versions
Detects providers without pinned versions, which can cause unexpected upgrades.

```hcl
# ❌ Bad - Missing version constraint
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

# ✅ Good - Version pinned
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

### 🔐 Hardcoded Credentials
Scans for hardcoded secrets that should never be committed.

Detects:
- AWS Access Keys (AKIA...)
- AWS Secret Keys
- Passwords in plaintext
- API Keys
- Generic secrets

**Fix**: Use environment variables or secret management systems:
```hcl
# ❌ Bad
resource "aws_instance" "example" {
  access_key = "AKIAIOSFODNN7EXAMPLE"
}

# ✅ Good
resource "aws_instance" "example" {
  access_key = var.aws_access_key
}
```

## Output Format

```
=== terraform Audit Results ===

SEVERITY  CHECK NAME             RESOURCE             MESSAGE
--------  ----------             --------             -------
HIGH      terraform-validate     ./infra              Terraform configuration is invalid
MEDIUM    provider-version       ./main.tf            Provider version constraint not found
CRITICAL  hardcoded-credentials  ./secrets.tf         Hardcoded aws_access_key detected
```

## Exit Codes

- `0`: All checks passed
- `1`: LOW severity issues found
- `2`: MEDIUM severity issues found
- `3`: HIGH severity issues found
- `4`: CRITICAL severity issues found

## Requirements

- Go 1.21+ (for devopsctl)
- Terraform binary in PATH (for format and validation checks)
- `.tf` files in the target directory

## Common Issues

### "terraform: not found"
The terraform binary is not installed or not in your PATH.
- Install: https://www.terraform.io/downloads.html
- Or skip: Format and validate checks won't run, but credential scanning will

### "directory does not exist"
Ensure the directory path is correct:
```bash
devopsctl validate terraform --dir /absolute/path/to/terraform
# or
cd /path/to/terraform && devopsctl validate terraform
```

### False positives on credential detection
If you have strings that match credential patterns but aren't credentials:
- Move the string to a variable
- Use comments to mark safe patterns
- Use secret management for actual credentials

## Integration

Use in CI/CD pipelines:
```yaml
# GitHub Actions example
- name: Validate Terraform
  run: devopsctl validate terraform --dir ./infrastructure
```
