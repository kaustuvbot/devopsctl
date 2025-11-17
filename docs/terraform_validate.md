# Terraform Validate Module Specification

Command:
devopsctl validate terraform

---

## 🎯 Purpose

Validate Terraform code quality and security posture.

---

## 🏗 Code Quality Checks

### 1. terraform fmt check
Run: terraform fmt -check

---

### 2. terraform validate
Run: terraform validate

---

### 3. Hardcoded Credentials Detection
Search:
- access_key
- secret_key
- password

Severity: HIGH

---

### 4. Unpinned Provider Versions
Detect:
version = "~>"
If missing → MEDIUM

---

## 🔐 Security Scan

Optional:
- Wrap tfsec
- Wrap checkov

Parse results.
Normalize severity.

---

## 🧠 Implementation Rules

- Use os/exec safely
- Capture stdout/stderr
- Do not crash if terraform not installed
