# Security Policy

## Supported Versions

DigiEmu Core follows a security-first maintenance model for tagged releases.

| Version line | Supported |
|---|---|
| main / next release candidate | yes |
| latest tagged stable release | yes |
| older releases | no guarantee |

Only the latest stable release line and the current main branch are considered security-supported unless explicitly stated otherwise.

## Security Scope

Security in DigiEmu Core includes, but is not limited to:

- Deterministic integrity violations
- Snapshot tampering
- Replay inconsistencies
- Canonicalization mismatches
- Hidden mutable state affecting verification
- Nondeterministic behavior in core paths
- Unauthorized or unsafe write behavior in bundle/update flows
- Dependency-level risks that can affect integrity, verification, or auditability
- Unsafe API behavior that can corrupt state, bypass validation, or weaken provenance

## What Is Considered Critical

Issues are treated as high priority if they can:

- break deterministic replay
- produce different hashes for semantically identical state
- allow silent mutation of semantic state
- bypass audit evidence
- overwrite expected hashes unsafely
- introduce hidden state or nondeterministic behavior in core logic
- compromise provenance, trust, or verification guarantees

## Reporting a Vulnerability

Please report suspected vulnerabilities privately and do not disclose them publicly before triage.

Include:

- affected version or commit
- impacted file(s) or package(s)
- reproduction steps
- expected behavior
- actual behavior
- whether the issue affects determinism, integrity, auditability, availability, or confidentiality
- proof-of-concept only if safe and minimal

Preferred report format:

1. Summary
2. Impact
3. Reproduction steps
4. Affected components
5. Suggested fix or mitigation
6. Whether public disclosure has occurred

## Response Goals

Target response goals for privately reported issues:

- initial acknowledgement: as soon as possible
- triage: best effort
- fix plan: after reproduction and impact confirmation
- disclosure: after fix or mitigation is available

These are goals, not contractual SLAs.

## Disclosure Policy

Please do not publish working exploit details before a fix or mitigation exists.

Coordinated disclosure is preferred.

When a vulnerability is confirmed, maintainers may publish:

- affected versions
- severity
- impact summary
- mitigation
- patch status
- upgrade path

## Security Design Principles

DigiEmu Core prioritizes:

- deterministic behavior over convenience
- explicit state over implicit state
- canonical serialization over serializer defaults
- append-only audit evidence
- fail-closed behavior for integrity-sensitive operations
- explicit documentation of accepted exceptions

## Out of Scope

Unless explicitly tied to a real integrity or security impact, the following are generally out of scope:

- stylistic code issues
- non-core CLI ergonomics
- hypothetical attacks without a reproducible path
- performance-only concerns without integrity impact
- warnings already documented as accepted deterministic exceptions

## Hardening Expectations

The project aims to maintain:

- threat model documentation
- dependency review
- documented deterministic exceptions
- test coverage for canonical replay and verification invariants
- controlled write policies for expected hash updates
- clear release and versioning rules

## No Warranty

Security review reduces risk but does not guarantee absence of defects.
DigiEmu Core should be deployed with operational controls, change review, and environment hardening.