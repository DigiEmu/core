# SECURITY_POLICY.md

## Security Contact

Security-related reports for DigiEmu Core should be sent privately to the project maintainer before public disclosure.

Recommended report content:
- affected version or commit
- reproduction steps
- expected behavior
- actual behavior
- security impact
- proof-of-concept if available
- whether the issue is publicly known

---

## Scope

This policy applies to:
- DigiEmu Core source code
- CLI verification and replay behavior
- canonicalization and hashing logic
- snapshot bundle loading
- audit integrity behavior
- public result contracts
- release artifacts maintained in this repository

This policy does **not** automatically cover:
- third-party deployment environments
- reverse proxies
- cloud infrastructure
- downstream applications integrating Core
- local workstation security
- secrets outside this repository

---

## Supported Security Goals

DigiEmu Core primarily protects:

- deterministic reconstruction
- tamper evidence
- auditability
- canonical serialization stability
- explicit failure on malformed or inconsistent state

DigiEmu Core does **not** claim to provide by itself:

- full application security
- authentication and authorization
- transport encryption
- secrets management
- malware resistance
- cloud tenancy isolation

---

## Reporting Guidelines

Please report privately if you find issues such as:

- nondeterministic behavior in verification or hashing
- inconsistent canonicalization
- unexpected bundle-loading order
- ability to overwrite expected hashes outside documented policy
- silent acceptance of malformed bundle state
- missing or bypassed audit behavior in strict-audit paths
- contract drift in result JSON that breaks verification consumers
- crashers or parser failures caused by malformed input
- evidence that semantic state can disappear unexpectedly during serialization

Please do **not** file public issues first for material security problems.

---

## Response Process

### 1. Acknowledgement
A valid report should be acknowledged as soon as reasonably possible.

### 2. Triage
The issue will be reviewed for:
- reproducibility
- impact on determinism
- impact on integrity
- impact on auditability
- impact on public verification contract
- exploitability in normal usage

### 3. Classification
Findings are typically classified into:
- critical
- high
- medium
- low
- informational

### 4. Remediation
If confirmed, the issue may be handled by:
- code fix
- documentation fix
- guard rule update
- explicit residual-risk documentation
- release-note disclosure
- versioned contract change

### 5. Disclosure
Public disclosure should happen only after:
- fix or mitigation is available, or
- project maintainer documents why risk is accepted

---

## Severity Guidance

### Critical
A flaw that breaks deterministic trust, allows undetected integrity compromise, or invalidates core audit evidence in a material way.

Examples:
- nondeterministic hashing in released verification path
- silent acceptance of tampered expected hash
- unbounded behavior causing systematic verification ambiguity

### High
A flaw that materially weakens verification reliability or audit integrity but still has partial mitigations.

Examples:
- unstable file ordering in bundle assembly
- contract-drifting result output in supported release path
- audit gaps in strict mutation flow

### Medium
A flaw that affects robustness, evidence clarity, or safe failure behavior.

Examples:
- malformed input causing confusing error semantics
- optional fields obscuring non-critical but relevant evidence
- environment-specific output drift with documented workaround absent

### Low
A minor hardening gap or edge case with limited direct security effect.

### Informational
A design observation, recommendation, or documentation improvement.

---

## Safe Harbor

Good-faith security research intended to improve DigiEmu Core is welcome when it avoids:
- destructive actions
- privacy violations
- service disruption beyond minimal validation
- unauthorized persistence or exfiltration
- misleading public claims before disclosure handling

---

## Security Update Policy

Security-relevant fixes should be documented through:
- commit history
- release notes
- updated docs where needed
- test coverage where appropriate

When a change affects public verification behavior or canonical output semantics, the project must also review:
- versioning impact
- compatibility expectations
- enterprise support implications

---

## Enterprise Distribution Requirement

Before marketing DigiEmu Core as enterprise-ready, the repository should include:

- published threat model
- published security policy
- dependency review
- documented support policy
- documented versioning guarantees
- documented deterministic exceptions / ignores
- passing build, test, and guard checks for the release candidate

---

## Known Limitations

Current repository health may still include:
- documented guard ignores for selected determinism findings
- accepted `omitempty` warnings in non-critical or contract-stable structures
- platform-specific output handling differences outside core execution logic

These limitations must be treated as explicit, documented risk decisions, not hidden assumptions.

---

## No Warranty by This Policy

This policy describes process and intent.  
It does not by itself create a warranty, SLA, or legal guarantee unless separately agreed in writing.