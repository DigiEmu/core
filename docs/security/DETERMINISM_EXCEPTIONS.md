
---

## `DETERMINISM_EXCEPTIONS.md`


# DETERMINISM_EXCEPTIONS.md

## Purpose

This document records known, intentional, or temporarily accepted determinism-related exceptions, ignores, and review findings for DigiEmu Core.

Its purpose is to prevent hidden risk from being confused with resolved risk.

---

## Principle

A documented exception is **not** the same as a fix.

It means one of the following:

- the finding is known and accepted temporarily
- the finding is operationally irrelevant to released trust semantics
- the finding is outside the authoritative deterministic path
- the finding requires later remediation but is not release-blocking under current policy

---

## Current Guard Ignore File

The repository currently uses:

- `.digiemu-guard-ignore.json`

This file must be treated as a policy artifact and reviewed before every release.

---

## Current Ignored Areas

## 1. `internal/kernel/usecases/snapshot_hash.go`
Ignored category/rule:
- determinism / range-review

### Reason
This path contains range-based constructs that triggered deterministic review warnings.  
The current release policy treats these findings as documented exceptions rather than silently unresolved issues.

### Required condition
If this file becomes part of a stronger enterprise trust claim, the ignores should be revisited and ideally eliminated through explicit ordering.

---

## 2. `internal/kernel/usecases/verify_audit.go`
Ignored category/rule:
- determinism / range-review

### Reason
Guard identified range/map-order concerns in audit-verification-related code.

### Current position
Accepted temporarily only because:
- behavior is known
- the ignore is explicit
- repository state is documented

### Required condition
Before claiming full audit-determinism hardening for enterprise contracts, this file should be reviewed line-by-line for key ordering guarantees.

---

## 3. `internal/kernel/usecases/export_unit_snapshot.go`
Ignored category/rule:
- determinism / range-review

### Reason
Guard flagged range-loop review in snapshot export logic.

### Current position
Accepted as a documented exception.

### Required condition
If export artifacts become contractual evidence outputs, ordering should be made explicit rather than merely reviewed.

---

## 4. `internal/tests/`
Ignored rule:
- range-review

### Reason
Test code may legitimately use patterns that are not release-critical for runtime trust semantics.

### Constraint
This exception must never be used to justify nondeterminism in production code.

---

## 5. `cmd/digiemu/`
Ignored category:
- time

### Reason
CLI/runtime utility behavior may reference time-related concerns outside the canonical deterministic state model.

### Constraint
Time usage must not leak into canonical hash computation or verification contract semantics.

---

## What Is Not Accepted as an Exception

The following are **not** acceptable to leave undocumented:

- map iteration affecting canonical hash
- unstable verification result contract
- bundle file ordering drift in authoritative verify path
- silent omission of integrity-relevant fields
- hidden environment-dependent replay changes
- undocumented write-policy weakening

If any of these occur, they are release blockers unless explicitly resolved or reclassified with strong justification.

---

## Exception Review Questions

For each ignored finding, ask:

1. Is the path on the authoritative verification or hashing route?
2. Can the finding affect byte-exact output?
3. Can the finding affect audit evidence?
4. Can two runs produce different trust conclusions?
5. Is this ignored because it is safe, or because it is unfinished?

If the honest answer is “unfinished and potentially trust-affecting,” it should not stay ignored indefinitely.

---

## Release Condition

A release may proceed with documented determinism exceptions only if:

- build passes
- tests pass
- guard passes with explicit ignore file
- ignored findings are listed here
- no ignored finding is known to alter released trust semantics

---

## Enterprise Positioning Guidance

For enterprise discussions, say:

- some determinism-related guard findings are currently documented as accepted exceptions
- the repository is transparent about these exceptions
- these are governance-visible risk decisions, not hidden defects

Do **not** say:
- “all determinism issues are fully solved”
- “guard clean means zero residual risk”
- “ignored findings are irrelevant” unless that has been explicitly demonstrated

---

## Planned Hardening Direction

Recommended future work:

- replace map iteration with sorted key traversal where deterministic output matters
- make ordering rules explicit in export and audit helpers
- reduce reliance on ignore policy over time
- distinguish authoritative deterministic paths from non-authoritative helper paths in code and docs

---

## Review Cadence

Review this document:
- before every deploy/tag
- after guard-rule changes
- after any new determinism warning
- before enterprise licensing conversations