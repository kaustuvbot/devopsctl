# Doctor Engine Specification

Command:
devopsctl doctor

---

## Purpose

Aggregate results from all audit and validation modules into a unified health report.
Act as the single entry point for running all checks.

---

## Module Orchestration

The doctor engine runs all registered modules:
- AWS Audit
- Docker Audit
- Terraform Validate
- Git Audit

Each module returns a list of CheckResult structs.
The doctor aggregates all results into a single report.

Modules that fail to execute (e.g., missing credentials, tool not installed)
must be handled gracefully with an error message, not a crash.

---

## Severity Scoring

Weighted severity model:
- LOW = 1
- MEDIUM = 2
- HIGH = 3
- CRITICAL = 4

Summary score = sum of all finding weights.
Exit code = highest severity found across all modules.

---

## Output Formats

### Table (default)
Colored terminal table grouped by module and severity.

### JSON (--json flag)
```json
{
  "modules": [
    {
      "module": "aws",
      "results": [...]
    }
  ],
  "summary": {
    "total_findings": 12,
    "critical": 2,
    "high": 4,
    "medium": 3,
    "low": 3,
    "score": 38
  }
}
```

### Markdown (--output report.md)
Generates a markdown report file with sections per module.

---

## Exit Codes

- 0 = No findings
- 1 = LOW severity findings present
- 2 = MEDIUM severity present
- 3 = HIGH severity present
- 4 = CRITICAL severity present

Exit code reflects the highest severity found.

---

## Code Rules

- Orchestrator must not import module internals directly.
- Use a common check interface for module registration.
- Each module registers itself with the doctor engine.
- No printing inside the doctor engine — use reporter module.
- Handle partial failures (one module fails, others still run).
