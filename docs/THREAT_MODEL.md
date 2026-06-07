# DigiEmu Core Threat Model

Version: 1.0  
Status: Draft for internal hardening and enterprise readiness  
Scope: DigiEmu Core repository, CLI verification flow, deterministic replay, bundle loading, audit verification, snapshot hashing

---

## 1. Purpose

This document identifies the main threats relevant to DigiEmu Core and defines the security posture required to operate the system as a deterministic, auditable knowledge infrastructure component.

DigiEmu Core is not a general-purpose web application. Its primary value lies in:

- deterministic replay
- canonical serialization
- stable snapshot hashing
- audit-ready verification output
- reproducible evidence generation

Accordingly, the main security objective is not only confidentiality or availability, but also **integrity of epistemic state and reproducibility of verification results**.

---

## 2. Security Objectives

The following objectives are mandatory:

1. **Deterministic integrity**
   The same valid bundle input must always reconstruct to the same canonical state and the same hash.

2. **Audit integrity**
   Verification and audit outputs must be reproducible and protected from nondeterministic drift.

3. **Canonical serialization safety**
   No semantic state required for verification may silently disappear due to encoding, omitted fields, or platform-specific formatting.

4. **Bundle parsing safety**
   Malformed, partial, oversized, or adversarial bundle input must fail safely and explicitly.

5. **Controlled mutation**
   Any workflow that modifies bundle content, such as `write-expected`, must be tightly constrained and explicitly classified.

6. **Operational predictability**
   CLI and library behavior must remain stable across supported platforms.

---

## 3. Assets

The following assets are security-relevant:

- snapshot bundles
- `snapshot.json`
- optional collection directories:
  - `units/`
  - `versions/`
  - `claims/`
  - `meaning/`
  - `uncertainty/`
- expected hash values
- deterministic verification reports
- audit event streams
- canonical JSON outputs
- replay results
- verification result contracts
- fixture bundles used as reference evidence

---

## 4. Trust Boundaries

### Trusted

- committed source code in the reviewed repository
- tagged release artifacts
- explicitly approved fixtures and example bundles
- canonical JSON implementation
- deterministic hashing implementation

### Conditionally trusted

- local filesystem
- CI environment
- CLI arguments
- example bundles in working tree
- generated reports written by automation

### Untrusted

- user-supplied bundles
- malformed JSON input
- partial filesystem contents
- path-manipulated bundle roots
- hostile or oversized collection payloads
- platform-specific output encodings introduced by shell redirection tools
- any map iteration order affecting public results

---

## 5. Threat Actors

### Accidental actor
A developer or operator unintentionally introduces drift through:
- field-tag changes
- `omitempty` usage
- map iteration
- output encoding conversion
- line-ending conversion
- fixture corruption

### Malicious local actor
A user with filesystem access attempts to:
- alter bundle contents
- inject malformed JSON
- exploit path assumptions
- replace expected outputs
- create misleading verification artifacts

### Supply-chain actor
A dependency or toolchain change introduces:
- changed JSON semantics
- changed filesystem behavior
- changed CLI output encoding
- altered deterministic ordering

### Integrator misuse
An integrator assumes behavior that Core does not guarantee, such as:
- treating warnings as approvals
- ignoring write-policy classifications
- using modified fixtures as normative reference

---

## 6. Main Threats

### T1. Nondeterministic ordering
**Description:** range over maps or unstable iteration changes output order, hash, trace, or audit evidence.

**Impact:** Critical integrity failure.

**Examples:**
- ranging over Go maps in snapshot hashing
- nondeterministic audit field ordering
- unstable collection traversal

**Mitigations:**
- never allow map iteration to affect hashed or public output
- sort keys explicitly before serialization or comparison
- sort filenames before bundle loading
- guard checks for range-review findings
- maintain explicit ignores only where behavior is documented and proven non-semantic

---

### T2. Silent semantic loss from `omitempty`
**Description:** fields disappear from serialized output, making two different semantic states appear identical.

**Impact:** High integrity risk.

**Examples:**
- verification result fields omitted unexpectedly
- bundle collections omitted instead of represented as empty lists
- domain objects losing explicit empty state

**Mitigations:**
- remove `omitempty` from hashed and contract-critical structures where absence must not be ambiguous
- initialize empty slices explicitly where contract requires `[]`
- document exceptions where omission is explicitly non-semantic
- review serialization warnings before release

---

### T3. Self-referential or unstable hash input
**Description:** hash calculation includes mutable or self-referential fields such as `expected_hash_v1`.

**Impact:** Critical correctness failure.

**Mitigations:**
- exclude `expected_hash_v1` from canonical hash scope
- keep canonical scope explicitly documented
- test byte-exact reproduction against committed fixtures

---

### T4. Unsafe write-back behavior
**Description:** `write-expected` overwrites state unexpectedly or accepts invalid hash material.

**Impact:** High integrity and audit risk.

**Mitigations:**
- only permit writes when placeholder policy allows it
- never overwrite existing expected values silently
- classify blocked writes with stable reasons
- reject invalid hash material
- return typed error categories

---

### T5. Malformed or hostile bundle content
**Description:** invalid JSON, oversized content, unexpected directory layout, BOM issues, or missing required files.

**Impact:** Medium to high.

**Mitigations:**
- BOM-safe reads
- explicit JSON decode failures
- bounded payload size in use cases
- strict path validation for bundle roots
- safe failure with deterministic error reporting

---

### T6. Output drift due to shell/tool encoding
**Description:** PowerShell, Windows redirection, or external tooling introduces UTF-16, BOM, CRLF/LF mismatches, or trailing newline drift.

**Impact:** High for byte-exact test fixtures and reproducibility claims.

**Mitigations:**
- treat CLI stdout as authoritative UTF-8 JSON
- read expected fixtures from Git blob when testing byte equality
- normalize only where test contract explicitly permits it
- avoid regenerating expected fixtures with shell tools unless encoding is controlled

---

### T7. Path confusion and unintended root selection
**Description:** bundle resolution loads from the wrong fixture or data root.

**Impact:** Medium.

**Mitigations:**
- explicit root trace recording
- stable bundle-root search order
- `prefer-data` opt-in only
- `--bundle` path validation requiring `snapshots/<ref>` structure

---

### T8. Audit verification drift
**Description:** audit comparison or hashing changes because of unstable field order or serialization ambiguity.

**Impact:** Critical for evidence trust.

**Mitigations:**
- canonicalize audit structures before hashing
- never depend on map iteration order
- keep audit comparison logic deterministic
- document ignored determinism cases if truly non-semantic

---

## 7. Threat Severity Summary

| Threat | Severity | Reason |
|---|---:|---|
| Nondeterministic ordering | Critical | Breaks core promise of reproducibility |
| Silent semantic loss | High | Can collapse distinct states |
| Self-referential hashing | Critical | Invalidates snapshot integrity |
| Unsafe write-back | High | Can falsify trusted bundle state |
| Malformed hostile input | Medium/High | May cause invalid or misleading results |
| Encoding/output drift | High | Breaks byte-exact evidence |
| Path confusion | Medium | Can verify wrong artifact |
| Audit drift | Critical | Breaks evidentiary trust |

---

## 8. Current Security Posture

At the current stage, DigiEmu Core is hardened primarily around:

- deterministic replay
- canonical scope control
- explicit write-policy handling
- guard-based detection of determinism and serialization issues
- reproducibility tests
- fixture-based byte-exact verification

Remaining enterprise hardening depends on:

- completion of security policy
- dependency review discipline
- documented support/versioning guarantees
- abuse-case documentation
- release checklist enforcement

---

## 9. Non-Goals

This threat model does not claim that DigiEmu Core currently provides:

- multi-tenant isolation
- network perimeter protection
- secret storage
- encryption-at-rest
- user authentication framework
- authorization policy engine
- anti-DDoS protections
- hosted SaaS operational guarantees

Those concerns belong to surrounding deployment and integration layers unless later added explicitly.

---

## 10. Required Controls Before Enterprise Licensing

Before positioning DigiEmu Core as enterprise-license ready, the following controls must be present and maintained:

- documented security policy
- dependency audit discipline
- release checklist
- explicit versioning policy
- deterministic test matrix
- abuse-case list
- supported-platform statement
- documented ignores for determinism reviews
- signed or otherwise controlled release process

---

## 11. Review Cadence

This threat model should be reviewed:

- before each tagged release
- whenever bundle format changes
- whenever canonical JSON or hash scope changes
- whenever write-policy semantics change
- whenever new API or transport surfaces are introduced

---

## 12. Approval Rule

No release should be described as enterprise-ready unless:

- build passes
- full test suite passes
- byte-exact repro tests pass
- determinism guard findings are either zero or explicitly documented and approved
- this threat model remains materially accurate
---

## Core 2.0 Verification Boundary Addendum

DigiEmu Core 2.0 is limited to deterministic decision-state reconstruction and verification.

It detects modified snapshots, hash mismatches, schema violations, non-canonical serialization problems, replay mismatches, malformed verification inputs, and incompatible schema or core versions.

It does not prove that an AI decision was ethically correct, that model reasoning was truthful, that a human operator acted responsibly, that an agent identity is trusted, that an action was authorized, that a deployed system is legally compliant, or that a model or agent is certified.

Agent identity, trust certification, authorization, and action attestation are outside DigiEmu Core and may be handled by complementary systems such as TBN.
