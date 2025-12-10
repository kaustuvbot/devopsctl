# Terraform Validate Module Specification

Command:
```
devopsctl validate terraform [--dir <path>]
```

---

## 🎯 Purpose

Validate Terraform code quality, configuration validity, and security posture by running multiple automated checks.

---

## 🏗️ Implementation Architecture

```
Runner (orchestrates all checks)
├── CheckFormat (terraform fmt -check)
├── CheckValidate (terraform validate)
├── CheckProviderVersions (unpinned provider detection)
└── CheckCredentials (hardcoded credential detection)
```

---

## ✅ Check Details

### 1. Format Check (terraform-fmt)
**Command**: `terraform fmt -check -recursive`
**Severity**: MEDIUM
**When triggered**: When terraform files don't match canonical formatting
**Recommendation**: Run `terraform fmt` to auto-fix

### 2. Validation Check (terraform-validate)
**Command**: `terraform validate`
**Severity**: HIGH
**When triggered**: When HCL syntax is invalid or configuration has errors
**Recommendation**: Fix terraform configuration errors per error message

### 3. Provider Version Check (provider-version)
**Type**: Pattern matching in HCL files
**Severity**: MEDIUM
**Pattern**: Detects `required_providers` block without `version` constraint
**Recommendation**: Add version constraints to provider configuration

### 4. Credentials Check (hardcoded-credentials)
**Type**: Regex pattern matching
**Severity**: CRITICAL
**Patterns detected**:
- AWS Access Keys: `AKIA[0-9A-Z]{16}`
- AWS Secret Keys: `[A-Za-z0-9/+=]{40}`
- Passwords: `password\s*=\s*"[^"]+"`
- API Keys: `api_key\s*=\s*"[^"]+"`
- Secrets: `secret\s*=\s*"[^"]+"`
**Recommendation**: Use environment variables or secret management (Vault, AWS Secrets Manager)

---

## 🔄 Check Execution Flow

1. Recursively find all `.tf` files in working directory
2. For format & validate: run terraform binary
3. For provider versions: parse HCL syntax, extract required_providers blocks
4. For credentials: scan file content with regex patterns
5. Aggregate results with severity levels
6. Return structured CheckResult array

---

## ⚠️ Error Handling

- If terraform binary not found: gracefully skip exec-based checks
- If directory doesn't exist: return error
- If file read fails: skip that file, continue checking others
- Invalid HCL files: skip, continue with other files

---

## 🧠 Implementation Rules

- Use `os/exec` safely with proper error handling
- Do not crash if terraform binary is not installed
- Return empty result arrays for missing files, not errors
- All checks return ([]CheckResult, error) signature
- Use severity.Level type for standardization
