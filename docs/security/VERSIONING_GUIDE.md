# DigiEmu Core Versioning Guarantees

Version: 1.0
Status: Draft

## Purpose

This document defines what DigiEmu Core intends to keep stable across versions, and what may change.

Because DigiEmu Core makes determinism and verification claims, versioning must protect not only APIs but also evidence semantics.

---

## Versioning Model

DigiEmu Core follows semantic intent, with additional caution for deterministic and evidence-related behavior.

Version classes:

- MAJOR: breaking contract or semantic behavior changes
- MINOR: backward-compatible capability additions
- PATCH: backward-compatible fixes, hardening, or documentation-aligned corrections

---

## Stability Priorities

The highest stability priority applies to:

- deterministic verification behavior
- canonical hashing semantics
- public verification result fields
- bundle interpretation rules
- write-policy semantics
- replay evidence semantics
- documented release-gate behavior

These must not drift casually.

---

## Guaranteed Stable Within a Major Version

Within a major version, DigiEmu Core aims to preserve:

### 1. Verification result contract
The meaning of core verification fields should remain stable, including:

- `ok`
- `ref`
- `expected`
- `got`
- `hash_alg`
- `canonical_scope`
- `trace`
- `errors`
- `message`

If write-policy fields are part of the active contract, they are also expected to remain stable.

### 2. Deterministic replay semantics
The same supported bundle input should reconstruct the same semantic replay state.

### 3. Canonical hashing semantics
The meaning of the declared hashing algorithm and canonical scope must not change silently.

### 4. Exit-behavior meaning
Verification failure may still emit valid structured stdout. This behavior must not be changed silently if documented as part of the contract.

### 5. Release artifact interpretation
Reference fixtures and normative verification artifacts must not be redefined silently.

---

## Changes That Require Explicit Major-Version Treatment

The following should be treated as major changes unless a very strong compatibility case exists:

- changing hashed scope semantics
- changing canonical serialization rules
- removing or renaming public verification fields
- changing meaning of `ok`, `expected`, or `got`
- changing trace semantics in a way that breaks documented consumers
- changing write-policy classification strings
- changing bundle-loading rules in a way that affects evidence semantics

---

## Changes Allowed in Minor Versions

Minor versions may add:

- new non-breaking CLI capabilities
- new optional fields that do not alter meaning of existing fields
- new documented commands
- new validation checks that reject previously invalid input
- new tests, docs, and hardening behavior
- internal refactors with no contract drift

Minor versions should not introduce silent output drift in normative examples without explicit documentation.

---

## Changes Allowed in Patch Versions

Patch versions may include:

- bug fixes
- security fixes
- deterministic ordering fixes
- fixture handling fixes
- encoding robustness improvements
- documentation corrections
- test harness corrections

A patch may tighten behavior if the previous behavior was clearly erroneous, unsafe, or inconsistent with documentation.

---

## Non-Guaranteed Areas

The following are not guaranteed stable unless explicitly documented:

- internal package structure
- internal helper names
- internal test helper behavior
- draft or experimental files
- unexported types
- undocumented CLI edge cases
- ad hoc fixture-writing workflows

---

## Normative Artifact Rules

Normative fixtures used for reproducibility are part of the effective contract when tests compare against them byte-exactly.

Therefore:

- fixture rewrites must be deliberate
- newline or encoding changes are not “cosmetic” when fixture bytes are normative
- byte-exact repro drift must be reviewed like contract drift

---

## Ignore File Governance

The existence of a guard ignore file does not weaken versioning guarantees.

Any ignored determinism or serialization case must be:

- narrow
- justified
- documented
- non-contradictory to published guarantees

Broad ignores must not be used to simulate compatibility.

---

## Compatibility Statement

Consumers should assume the following rule:

If you rely on documented DigiEmu Core verification behavior within the same major version, that behavior should remain stable unless explicitly announced otherwise.

If a change affects deterministic evidence, it must be documented as a contract-relevant change.

---

## Release Notes Requirement

Every release should indicate whether it contains:

- no contract changes
- additive contract changes
- deterministic behavior fixes
- fixture or evidence artifact changes
- major semantic changes

This is especially important for enterprise users and audit-sensitive consumers.

---

## Recommended Tagging Discipline

For strong operational clarity, each release should clearly identify:

- code version
- contract version if separately tracked
- fixture version if normative artifacts changed
- schema version where applicable

---

## Final Rule

If a reasonable integrator would get a different verification meaning from the same documented workflow after upgrading, that change must be treated as version-significant and explicitly documented.