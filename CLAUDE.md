# CLAUDE.md
Project: devopsctl
Type: Modular DevOps Infrastructure Hygiene Toolkit
Language: Go
Primary OS Target: Ubuntu Linux

---

## 🎯 Purpose
devopsctl is a lightweight CLI toolkit for infrastructure hygiene and DevOps validation across:
- AWS
- Docker
- Terraform
- Git
- Doctor (aggregator)

This file defines development discipline.
Functional behavior is defined in PROJECT_SPEC.md and /claude_docs.

---

## 🐧 Installation Requirement
The project must:
- Be installable via `go install`
- Use Go modules (go.mod)
- Define main entrypoint in cmd/devopsctl/main.go
- Work on Ubuntu
- Support Go 1.21+
- Produce static binary

---

## 🧱 Architecture Rules
1. Modular structure
2. No monolithic logic
3. CLI separated from business logic
4. No hardcoded thresholds
5. Config-driven via .devopsctl.yaml
6. Strongly typed
7. No fmt.Println inside checks
8. Reporter handles output formatting
9. Command execution must be safe

---

## 📂 Structure
```
devopsctl/
├── cmd/devopsctl/main.go
├── internal/
│   ├── cli/
│   ├── config/
│   ├── reporter/
│   ├── aws/
│   ├── docker/
│   ├── terraform/
│   ├── git/
│   └── doctor/
├── claude_docs/
│   └── commit_rule.md
|   └── other feature listing md files.md
├── tests/
├── go.mod
└── PROJECT_SPEC.md
```

---

## 🧠 Development Discipline
Claude must:
- Build incrementally
- Follow 10-20 commit batch discipline (see claude_docs/commit_rule.md)
- Never generate full project at once
- Keep commits realistic
- Improve tests gradually
- Keep claude_docs and code synchronized

---

## 📊 Output Standard
Each check must return:
```go
type CheckResult struct {
    CheckName      string
    Severity       string // LOW|MEDIUM|HIGH|CRITICAL
    ResourceID     string
    Message        string
    Recommendation string
}
```

Exit code must reflect highest severity found.

---

## 📚 Documentation Authority
- PROJECT_SPEC.md defines feature scope
- /claude_docs defines module contracts
- claude_docs/commit_rule.md defines commit strategy
- Documentation-first development is mandatory

---

## 🚫 Non-Goals
- No GUI
- No SaaS dashboard
- No complex async patterns
- No heavy frameworks

---

## 🔄 Commit Strategy
- **Timeline**: Oct 1, 2025 → Feb 15, 2026
- **Total commits**: ~150 commits
- **Pattern**: Natural development rhythm
- **See**: claude_docs/commit_rule.md