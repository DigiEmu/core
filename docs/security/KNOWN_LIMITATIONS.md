# DigiEmu Core Known Limitations

Version: 1.0
Status: Draft

## Purpose

This document lists known limitations of DigiEmu Core in its current state.

The goal is clarity, not weakness. A trustworthy infrastructure component states both what it guarantees and what it does not yet guarantee.

---

## 1. Current Positioning Limitation

DigiEmu Core is already strong in deterministic verification and replay discipline, but it is not automatically equivalent to full enterprise certification, regulatory certification, or managed production support.

Additional process, documentation, and operational controls may still be required depending on deployment context.

---

## 2. Platform / Shell Limitations

### Windows shell encoding sensitivity
Normative JSON fixtures can be affected by shell and editor behavior on Windows, especially:

- PowerShell redirection
- UTF-16 output
- BOM insertion
- CRLF/LF conversion

This does not invalidate Core logic, but it can invalidate byte-exact artifact comparisons if unmanaged.

### Tooling-dependent fixture generation
Not all shell commands produce byte-identical fixture files. Controlled generation procedures are required for normative artifacts.

---

## 3. Serialization Limitations

### `omitempty` review still exists as a governance concern
Some structures may still require careful review to ensure semantic state does not disappear unintentionally during serialization.

A warning-free guard run is strong evidence of control, but future edits could reintroduce ambiguity unless review discipline remains active.

### Empty-vs-omitted semantics
Some data models can be sensitive to the distinction between:

- missing field
- empty object
- empty array
- zero-value scalar

This must be handled deliberately in contract-critical structures.

---

## 4. Determinism Limitations

### Determinism depends on coding discipline
Go allows constructs such as map iteration that are nondeterministic by default. DigiEmu Core mitigates this through code review, tests, and guard rules, but the language does not prevent accidental reintroduction.

### Ignore-file risk
A narrow ignore file is currently used for accepted review cases. This is manageable, but it creates a governance requirement: ignores must remain narrow and documented.

---

## 5. Fixture / Repro Limitations

### Byte-exact repro tests are fragile by nature
They are intentionally strict, but this means they can fail because of:

- newline differences
- encoding changes
- BOM differences
- altered fixture-writing method

This is a feature for integrity, not a flaw, but it increases maintenance sensitivity.

### Expected report artifacts are normative
Reference verification reports are not ordinary convenience files. They behave more like contract fixtures and must be treated accordingly.

---

## 6. CLI Limitations

### Non-zero exit with valid JSON may confuse wrappers
The verify command can emit valid JSON to stdout even when verification fails and the process exits non-zero. Some wrappers wrongly treat non-zero exit as absence of valid result.

Consumers must follow the documented contract rather than simplistic shell assumptions.

### Shell redirection is not a trust mechanism
CLI output can be redirected, but redirection alone is not a controlled evidence-generation workflow.

---

## 7. Audit Limitations

### Audit strength depends on canonical interpretation
Audit data is only as strong as the canonical rules applied to it. If an integrator consumes audit events outside the documented canonical semantics, reproducibility may degrade.

### Optional audit data may still require tightening over time
As the system matures, some audit-adjacent fields may need stronger explicitness to avoid ambiguity in enterprise settings.

---

## 8. Support / Operations Limitations

### No implicit long-term support guarantee
Unless separately documented, current support expectations should be understood as release-based, not unlimited LTS.

### No implicit managed-service guarantees
Core is a software component, not by default a hosted managed compliance platform.

---

## 9. Security Limitations

### Dependency and platform risk still exist
Even with strong core logic, the system depends on:

- Go toolchain behavior
- operating system behavior
- filesystem semantics
- repository hygiene
- release discipline

Security is therefore not only a code property but also an operational property.

### Misuse by operators remains possible
Core can be used incorrectly even when the code is correct, especially around:

- fixture regeneration
- wrong bundle comparison
- encoding drift
- unsupported shell workflows

---

## 10. Enterprise Readiness Limitation

DigiEmu Core may be technically strong enough for enterprise licensing discussions, but “enterprise-ready” in the strongest sense usually also requires:

- threat model
- security policy
- dependency audit
- support policy
- versioning guarantees
- pre-deploy checklist
- abuse-case documentation
- documented release process

The software and the governance package must mature together.

---

## 11. Recommended Operational Response

These limitations are best managed by:

- strict release gating
- narrow ignore governance
- controlled fixture generation
- documented Windows-safe procedures
- continued serialization review
- explicit support and versioning policies
- periodic dependency and threat review

---

## 12. Final Statement

Known limitations do not negate DigiEmu Core’s value. They define the boundary between:

- what is already strong and credible now
- what must still be operationalized for higher-assurance enterprise deployment