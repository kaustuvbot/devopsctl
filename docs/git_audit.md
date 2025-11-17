# Git Audit Module Specification

Command:
devopsctl audit git

---

## 🎯 Purpose

Assess repository hygiene and maintenance quality.

---

## 📦 Repo Size Check
Use:
git count-objects -vH

If > threshold (e.g., 500MB) → MEDIUM

---

## 🌿 Stale Branches
Logic:
- List branches
- Get last commit date
- If >90 days old → LOW

---

## 🗂 Large Files
Scan for files >50MB

Severity: MEDIUM

---

## 🔒 Missing Branch Protection (Future via GitHub API)

If default branch unprotected → HIGH

---

## 🧠 Code Rules

- Use os/exec or go-git library
- Support local-only mode
- Future: GitHub API integration
