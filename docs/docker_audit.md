# Docker Audit Module Specification

Command:
devopsctl audit docker --file Dockerfile

---

## 🎯 Purpose

Analyze Dockerfile and optionally container image for security and optimization.

---

## 🐳 Dockerfile Static Checks

### 1. Latest Tag Usage
If "FROM ubuntu:latest" → flag
Severity: MEDIUM

---

### 2. Running as Root
If no USER directive → HIGH

---

### 3. No Healthcheck
If missing HEALTHCHECK → LOW

---

### 4. No Multi-stage Build
If single FROM → LOW

---

### 5. Exposed All Ports
If EXPOSE 0-65535 or risky ports → MEDIUM

---

## 🔎 Image Scan (Optional Flag)

devopsctl audit docker --image myimage

Wrap Trivy.
Parse JSON.
Extract CRITICAL & HIGH vulns.

---

## 🧠 Code Rules

- Dockerfile parser must not use naive string matching only.
- Use regex or lightweight parsing logic in Go.
- Return structured findings.
- Use os/exec for optional Trivy integration.
