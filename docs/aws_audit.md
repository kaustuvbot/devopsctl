# AWS Audit Module Specification

Command:
devopsctl audit aws

---

## 🎯 Purpose

Perform lightweight infrastructure hygiene checks across AWS account.

Use AWS SDK for Go v2.
Respect AWS profiles.
Support multi-region scanning.

---

## 🔐 IAM Checks

### 1. Users Without MFA
Logic:
- List IAM users
- Check MFADevices
- If empty → flag

Severity: HIGH
Recommendation: Enforce MFA

---

### 2. Old Access Keys
Logic:
- List access keys
- Calculate age in days
- If > threshold (default 90 days) → flag

Severity:
  90-120 days → MEDIUM
  >120 days → HIGH

---

### 3. Administrator Access Users
Logic:
- Detect users/groups attached to AdministratorAccess policy

Severity: CRITICAL

---

## ☁️ S3 Checks

### 4. Public Buckets
Logic:
- Check bucket policy
- Check public access block
- If public → flag

Severity: CRITICAL

---

### 5. No Encryption
Logic:
- Check server-side encryption config
- If none → flag

Severity: HIGH

---

### 6. No Versioning
Logic:
- Check versioning status
- If disabled → LOW

---

## 🌐 Security Group Checks

### 7. 0.0.0.0/0 SSH Open
Logic:
- Ingress rule
- Port 22
- CIDR 0.0.0.0/0

Severity: CRITICAL

---

### 8. 0.0.0.0/0 All Ports
Severity: CRITICAL

---

## 💾 EBS Checks

### 9. Unencrypted Volumes
Severity: HIGH

### 10. Unattached Volumes
Severity: LOW

---

## 📊 Output Format

Return structured result list.
No printing inside check functions.
Reporting handled by reporter module.

---

## 🧠 Code Expectations

- Separate service-specific files:
    iam_checks.go
    s3_checks.go
    ec2_checks.go

- Shared AWS client helper
- Handle permission errors gracefully
