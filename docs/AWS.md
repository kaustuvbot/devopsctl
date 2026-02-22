# AWS Audit

The `audit aws` command scans your AWS account and reports security risks and cost issues across IAM, S3, EC2 security groups, and EBS volumes.

## Prerequisites

Before running the AWS audit, you need:

1. **An AWS account** with programmatic access enabled
2. **AWS CLI installed** (optional but recommended for credential setup)
3. **AWS credentials configured** — devopsctl uses the standard AWS credential chain

### Setting up AWS credentials

**Option A: AWS CLI (recommended for local use)**
```bash
# Install the AWS CLI first, then:
aws configure

# You will be prompted for:
# AWS Access Key ID: AKIA...
# AWS Secret Access Key: ...
# Default region name: us-east-1
# Default output format: json
```

**Option B: Environment variables**
```bash
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export AWS_DEFAULT_REGION=us-east-1
```

**Option C: AWS credentials file** (`~/.aws/credentials`)
```ini
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY

[staging]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
```

### Minimum IAM permissions

The audit user needs read-only permissions. Attach the following policy to the IAM user or role running devopsctl:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DevopsctlAuditReadOnly",
      "Effect": "Allow",
      "Action": [
        "iam:ListUsers",
        "iam:ListMFADevices",
        "iam:ListAccessKeys",
        "iam:ListAttachedUserPolicies",
        "iam:ListGroupsForUser",
        "iam:ListAttachedGroupPolicies",
        "s3:ListAllMyBuckets",
        "s3:GetBucketAcl",
        "s3:GetBucketPublicAccessBlock",
        "s3:GetEncryptionConfiguration",
        "s3:GetBucketVersioning",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeVolumes"
      ],
      "Resource": "*"
    }
  ]
}
```

> **Note**: If the tool lacks certain permissions, it skips those checks gracefully and continues with the rest. You will not see errors — just fewer results.

---

## Installation

```bash
go install github.com/kaustuvprajapati/devopsctl/cmd/devopsctl@latest
```

Requires Go 1.21+. The binary is statically linked and has no runtime dependencies.

---

## Quick Start

```bash
# Run the full AWS audit using your default AWS credentials
devopsctl audit aws

# Audit a specific region
# (set region in .devopsctl.yaml — see Configuration section)

# Show only critical and high severity findings
devopsctl audit aws --quiet

# Output as JSON
devopsctl audit aws --format json

# Output as Markdown and save to a file
devopsctl audit aws --format markdown --output aws-report.md
```

---

## Checks Performed

The audit runs 9 checks across 4 AWS services. Each check is independent — a failure in one does not stop the others.

---

### IAM Checks

#### `iam-mfa-disabled` — Severity: HIGH

**What it checks**: Whether each IAM user has at least one MFA device registered.

**Why it matters**: Passwords alone are not enough. If a user's password is leaked, an attacker can log in to your AWS console without MFA. Enabling MFA adds a second factor (phone app, hardware key) that the attacker would also need.

**Example finding**:
```
HIGH    iam-mfa-disabled    alice    IAM user "alice" has no MFA device enabled
```

**How to fix**:
1. Go to AWS Console → IAM → Users → select the user
2. Click the **Security credentials** tab
3. Under **Multi-factor authentication (MFA)**, click **Assign MFA device**
4. Follow the setup wizard for a virtual MFA app (e.g., Google Authenticator, Authy)

---

#### `iam-old-access-key` — Severity: MEDIUM / HIGH

**What it checks**: Whether any active access keys are older than the configured threshold (default: 90 days). Keys older than 90 days are MEDIUM; keys older than 120 days are HIGH.

**Why it matters**: Long-lived access keys are a major security risk. If a key is accidentally committed to a public repository or leaked in logs, it stays valid indefinitely unless rotated. Regular rotation limits the damage window.

**Example finding**:
```
MEDIUM    iam-old-access-key    alice    Access key AKIA... for user "alice" is 95 days old (threshold: 90)
HIGH      iam-old-access-key    bob      Access key AKIA... for user "bob" is 130 days old
```

**How to fix**:
1. Create a new access key for the user
2. Update your applications/scripts to use the new key
3. Delete the old key

**Config**: Adjust the threshold in `.devopsctl.yaml`:
```yaml
aws:
  key_age_days: 90  # flag keys older than this many days
```

---

#### `iam-admin-access` — Severity: CRITICAL

**What it checks**: Whether any IAM user has the `AdministratorAccess` AWS managed policy attached — either directly or via a group.

**Why it matters**: Admin access gives full control over your entire AWS account. If this user's credentials are compromised, an attacker can create new users, delete resources, exfiltrate data, or run up costs. The principle of least privilege means users should only have the permissions they actually need.

**Example finding**:
```
CRITICAL    iam-admin-access    charlie    IAM user "charlie" has AdministratorAccess (direct policy)
CRITICAL    iam-admin-access    dave       IAM user "dave" has AdministratorAccess (via group "admins")
```

**How to fix**:
- Remove `AdministratorAccess` from users and replace with scoped policies
- Keep admin access only for break-glass or root account scenarios
- Use AWS Organizations SCPs to limit blast radius

---

### S3 Checks

#### `s3-public-bucket` — Severity: CRITICAL

**What it checks**: Whether any S3 bucket is publicly accessible. The check looks at both:
- The bucket's **Public Access Block** settings (all 4 flags must be enabled)
- The bucket **ACL** for grants to `AllUsers` or `AuthenticatedUsers`

**Why it matters**: A public S3 bucket means anyone on the internet can list or download your files. Many data breaches have happened due to accidentally public S3 buckets containing customer data, database backups, or internal documents.

**Example finding**:
```
CRITICAL    s3-public-bucket    my-backup-bucket    S3 bucket "my-backup-bucket" is publicly accessible
```

**How to fix**:
1. Go to S3 Console → select the bucket → **Permissions** tab
2. Under **Block public access**, click Edit and enable all 4 options
3. Check **Bucket ACL** and remove any grants to `All users` or `Authenticated users`
4. Apply the same at the account level: S3 → **Block Public Access (account settings)**

---

#### `s3-no-encryption` — Severity: HIGH

**What it checks**: Whether each S3 bucket has server-side encryption configured.

**Why it matters**: Without encryption, data stored in S3 is at rest in plaintext. Encryption ensures that even if someone gains unauthorized access to the underlying storage, the data is unreadable without the encryption keys.

**Example finding**:
```
HIGH    s3-no-encryption    logs-bucket    S3 bucket "logs-bucket" has no server-side encryption configured
```

**How to fix**:
1. Go to S3 Console → select the bucket → **Properties** tab
2. Under **Default encryption**, click Edit
3. Select **SSE-S3** (Amazon managed keys, simple) or **SSE-KMS** (customer managed keys, more control)
4. Save changes

---

#### `s3-versioning-disabled` — Severity: LOW

**What it checks**: Whether versioning is enabled on each S3 bucket.

**Why it matters**: Without versioning, if a file is accidentally deleted or overwritten, it is gone permanently. Versioning keeps a history of all object versions, allowing point-in-time recovery.

**Example finding**:
```
LOW    s3-versioning-disabled    assets-bucket    S3 bucket "assets-bucket" does not have versioning enabled
```

**How to fix**:
1. Go to S3 Console → select the bucket → **Properties** tab
2. Under **Bucket Versioning**, click Edit
3. Enable versioning

---

### EC2 / Security Group Checks

#### `sg-all-ports-open` — Severity: CRITICAL

**What it checks**: Whether any security group has an inbound rule allowing all traffic (`protocol: -1`) from anywhere on the internet (`0.0.0.0/0` or `::/0`).

**Why it matters**: This is the most permissive rule possible. It allows any type of traffic (TCP, UDP, ICMP, etc.) on any port from any IP address. Attackers actively scan the internet for such open instances.

**Example finding**:
```
CRITICAL    sg-all-ports-open    sg-0a1b2c3d    Security group "sg-0a1b2c3d" (web-servers) allows all traffic from 0.0.0.0/0
```

**How to fix**:
- Remove the catch-all rule
- Add specific rules for only the ports and protocols your application needs (e.g., TCP 443 for HTTPS)
- Restrict source IP ranges to known CIDRs where possible

---

#### `sg-ssh-open` — Severity: CRITICAL

**What it checks**: Whether any security group allows inbound SSH (port 22) from `0.0.0.0/0` (the entire internet).

**Why it matters**: Port 22 is the default SSH port. Bots continuously scan the internet for open port 22 and attempt brute-force login attacks. Exposing SSH to the internet is a common attack vector.

**Example finding**:
```
CRITICAL    sg-ssh-open    sg-0a1b2c3d    Security group "sg-0a1b2c3d" allows SSH from 0.0.0.0/0
```

**How to fix**:
- Restrict SSH to your specific office or VPN IP range: `1.2.3.4/32`
- Better: Use **AWS Systems Manager Session Manager** instead of SSH — no open ports needed
- Or: Use a bastion host that only your team can reach

---

### EBS Checks

#### `ebs-unencrypted` — Severity: HIGH

**What it checks**: Whether each EBS volume has encryption enabled.

**Why it matters**: EBS volumes store the persistent storage for your EC2 instances (OS, databases, application data). Without encryption, data is stored in plaintext on physical disks. EBS encryption is a one-click setting with no performance penalty.

**Example finding**:
```
HIGH    ebs-unencrypted    vol-0a1b2c3d4e5f    EBS volume "vol-0a1b2c3d4e5f" is not encrypted
```

**How to fix**:
- For new volumes: Enable **Encryption by default** at the account level (EC2 Console → Settings → Data protection)
- For existing volumes: Create an encrypted snapshot, then restore a new encrypted volume from it

---

#### `ebs-unattached` — Severity: LOW

**What it checks**: Whether any EBS volumes exist in `available` state — meaning they are not attached to any EC2 instance.

**Why it matters**: Unattached EBS volumes still incur storage costs even when no instance is using them. These are often forgotten volumes from terminated instances.

**Example finding**:
```
LOW    ebs-unattached    vol-0a1b2c3d4e5f    EBS volume "vol-0a1b2c3d4e5f" is not attached to any instance
```

**How to fix**:
- Verify the volume is not needed (check its name/tags and creation date)
- Take a final snapshot if you might need the data later
- Delete the volume in the EC2 Console → Volumes

---

## Configuration

Create a `.devopsctl.yaml` file in your project directory to customize the audit:

```yaml
aws:
  enabled: true          # Set to false to skip AWS in `devopsctl doctor` runs
  region: us-east-1     # AWS region to scan (default: us-east-1)
  profile: default      # AWS CLI profile to use (empty = use default credential chain)
  key_age_days: 90      # Threshold for flagging old access keys (default: 90)
```

### Config file locations

devopsctl searches for the config file in this order:
1. Path specified via `--config` flag
2. `.devopsctl.yaml` in the current directory
3. `.devopsctl.yml` in the current directory
4. Built-in defaults (if no file found)

### Using a named AWS profile

If you have multiple AWS accounts, configure a named profile:

```yaml
# .devopsctl.yaml
aws:
  region: eu-west-1
  profile: production
```

Then your `~/.aws/credentials` would have:
```ini
[production]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
```

### Suppressing specific checks

To ignore checks that don't apply to your setup:

```yaml
# .devopsctl.yaml
ignore:
  checks:
    - s3-versioning-disabled    # we use lifecycle policies instead
    - ebs-unattached            # managed by our cleanup automation
```

Check names to use in the ignore list:
- `iam-mfa-disabled`
- `iam-old-access-key`
- `iam-admin-access`
- `s3-public-bucket`
- `s3-no-encryption`
- `s3-versioning-disabled`
- `sg-all-ports-open`
- `sg-ssh-open`
- `ebs-unencrypted`
- `ebs-unattached`

---

## Output Formats

### Table (default)

Color-coded terminal output. CRITICAL findings are red, HIGH are yellow.

```
=== aws Audit Results ===

SEVERITY    CHECK NAME           RESOURCE            MESSAGE
--------    ----------           --------            -------
CRITICAL    s3-public-bucket     my-backup-bucket    S3 bucket "my-backup-bucket" is publicly accessible
CRITICAL    iam-admin-access     charlie             IAM user "charlie" has AdministratorAccess (direct policy)
HIGH        iam-mfa-disabled     alice               IAM user "alice" has no MFA device enabled
HIGH        ebs-unencrypted      vol-0a1b2c3d        EBS volume "vol-0a1b2c3d" is not encrypted
LOW         ebs-unattached       vol-0x9y8z7w        EBS volume "vol-0x9y8z7w" is not attached to any instance
```

### JSON

```bash
devopsctl audit aws --format json
```

```json
{
  "module": "aws",
  "results": [
    {
      "check_name": "s3-public-bucket",
      "severity": "CRITICAL",
      "resource_id": "my-backup-bucket",
      "message": "S3 bucket \"my-backup-bucket\" is publicly accessible",
      "recommendation": "Enable S3 Block Public Access settings for the bucket and account"
    },
    {
      "check_name": "iam-mfa-disabled",
      "severity": "HIGH",
      "resource_id": "alice",
      "message": "IAM user \"alice\" has no MFA device enabled",
      "recommendation": "Enable MFA for all IAM users"
    }
  ]
}
```

### Markdown

```bash
devopsctl audit aws --format markdown --output aws-report.md
```

Generates a Markdown table with a Recommendations section — useful for sharing findings with your team or including in pull requests.

---

## Exit Codes

The exit code reflects the highest severity finding:

| Code | Meaning |
|------|---------|
| `0` | All checks passed — no findings |
| `1` | LOW severity findings only |
| `2` | MEDIUM severity findings present |
| `3` | HIGH severity findings present |
| `4` | CRITICAL severity findings present |

Use exit codes in CI to fail pipelines on critical findings:
```bash
devopsctl audit aws
if [ $? -ge 4 ]; then
  echo "CRITICAL AWS security issues found. Blocking deploy."
  exit 1
fi
```

---

## Common Issues

### "No findings returned" — account looks clean

This can mean one of three things:
1. Your account genuinely has no issues (great!)
2. The audit user lacks permissions for some checks — those checks are silently skipped
3. The region in your config does not match where your resources live

To diagnose: run with `--format json` and check if `results` is empty, then verify your permissions match the [minimum IAM policy](#minimum-iam-permissions) above.

### "AccessDenied" errors

devopsctl handles permission errors gracefully — it skips the check and moves on. You will not see an error in the output; you will simply not see findings for that check.

To get full coverage, ensure your audit user has all permissions listed in the [Minimum IAM permissions](#minimum-iam-permissions) section.

### "NoCredentialProviders" — credentials not found

devopsctl cannot find AWS credentials. Set them up via one of the three methods in [Prerequisites](#prerequisites).

### Many findings — where to start?

Prioritize by severity:
1. **CRITICAL first**: Public S3 buckets, admin users, wide-open security groups — fix these immediately
2. **HIGH next**: Unencrypted storage, missing MFA, very old access keys
3. **LOW last**: Unattached volumes and disabled versioning are cost/resilience concerns, not immediate security risks

Use `--quiet` to see only CRITICAL and HIGH:
```bash
devopsctl audit aws --quiet
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: AWS Security Audit

on:
  schedule:
    - cron: '0 9 * * 1'  # Every Monday at 9am
  push:
    branches: [main]

jobs:
  aws-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install devopsctl
        run: go install github.com/kaustuvprajapati/devopsctl/cmd/devopsctl@latest

      - name: Run AWS audit
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_DEFAULT_REGION: us-east-1
        run: devopsctl audit aws --quiet

      - name: Save Markdown report
        if: always()
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_DEFAULT_REGION: us-east-1
        run: devopsctl audit aws --format markdown --output aws-audit-report.md

      - name: Upload report
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: aws-audit-report
          path: aws-audit-report.md
```

> **Tip**: Store AWS credentials in GitHub Secrets, not in your code. Create a dedicated IAM user with only the [minimum read-only permissions](#minimum-iam-permissions) for CI.
