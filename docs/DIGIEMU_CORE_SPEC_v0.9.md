# DigiEmu Core Specification v0.9

Status: Public Review Draft  
Version: 0.9  
Repository: DigiEmu/core  
Scope: Deterministic knowledge snapshot infrastructure for reproducible and auditable AI-related decision states.

---

## 1. Purpose

DigiEmu Core defines a deterministic infrastructure layer for capturing, canonicalizing and verifying knowledge or decision states in AI-assisted and complex digital systems.

The central verification claim is:

> Same input -> same reconstructed state -> same hash -> PASS / FAIL.

DigiEmu Core does not attempt to prove that an AI decision is correct, ethical or complete. It defines whether a specific state boundary can be reconstructed deterministically and verified independently.

---

## 2. Design Goals

DigiEmu Core is designed to provide:

- deterministic state capture
- canonical JSON serialization
- SHA-256 state identity
- clear inside-hash vs outside-hash separation
- replay-ready verification
- machine-readable PASS / FAIL results
- portable snapshot bundles
- audit-ready evidence for governance and technical review

---

## 3. Non-Goals

DigiEmu Core is not:

- a general AI safety framework
- a model evaluation benchmark
- a replacement for human governance
- a medical, legal or financial decision engine
- a blockchain
- a database
- a runtime permission system
- a guarantee that a decision is substantively correct

DigiEmu Core verifies reproducibility and integrity of deterministic state boundaries. It does not validate the moral, legal or domain-specific correctness of the decision itself.

---

## 4. Terminology

### Deterministic State

A state whose relevant fields can be reconstructed in a stable way and serialized into a canonical representation.

### Canonical Snapshot

A JSON document that represents the deterministic state boundary.

### State Boundary

The explicit separation between data that belongs inside the deterministic hash and metadata that remains outside the hash.

### Inside Hash

Fields that directly define the reproducible state and therefore affect the computed hash.

### Outside Hash

Fields that provide audit or operational context but must not affect deterministic state identity, such as timestamps, comments, runtime notes or environment metadata.

### Snapshot Hash

The SHA-256 digest computed over the canonical JSON representation of the deterministic snapshot.

### Replay Verification

The process of reconstructing a state, canonicalizing it, hashing it and comparing the result with an expected hash.

### Verification Result

A machine-readable PASS / FAIL result describing whether the reconstructed state matches the expected deterministic identity.

---

## 5. Core Model

DigiEmu Core follows this minimal model:

```text
Input
  -> Deterministic State Boundary
  -> Canonical Snapshot
  -> Canonical JSON
  -> SHA-256 Hash
  -> Replay Verification
  -> PASS / FAIL
```

A conforming DigiEmu Core implementation must be able to:

1. identify the deterministic state boundary
2. produce or consume a canonical snapshot
3. serialize the snapshot deterministically
4. compute the required hash
5. compare the computed hash with the expected hash
6. return a clear verification result

---

## 6. Snapshot Boundary

The snapshot boundary is the most important concept in DigiEmu Core.

Only deterministic facts that define the reproducible state belong inside the hash.

Examples of inside-hash fields:

- stable identifiers
- normalized state values
- declared meaning fields
- versioned deterministic units
- explicit uncertainty values when part of the state

Examples of outside-hash fields:

- timestamps
- user comments
- runtime environment notes
- UI labels
- log formatting
- human-readable audit notes
- non-deterministic execution metadata

This separation prevents accidental hash drift caused by operational metadata.

---

## 7. Canonical JSON v1

DigiEmu Core uses canonical JSON serialization to ensure that the same deterministic state produces the same byte representation before hashing.

A canonical JSON implementation must define stable handling of:

- object key ordering
- whitespace
- string encoding
- number representation
- null values
- arrays
- nested structures

The normative rules for DigiEmu canonical snapshot hashing are defined in:

- `docs/SNAPSHOT_HASH_v1.0.md`

---

## 8. Snapshot Hash v1

The snapshot hash is the SHA-256 digest of the canonical JSON representation of the deterministic state.

The expected result is represented as an uppercase hexadecimal string.

The snapshot hash is not a trust claim by itself. It is a deterministic identity claim.

A matching hash means:

> The reconstructed deterministic state is byte-equivalent to the expected canonical state.

A non-matching hash means:

> The reconstructed deterministic state differs from the expected canonical state or the canonicalization process is not conformant.

---

## 9. Bundle Layout v1

A DigiEmu bundle is a portable directory or archive containing the data required for independent verification.

A minimal bundle may contain:

```text
snapshot.json
expected_hash_v1.txt
verify_report.json
metadata.json
```

The bundle layout is defined in:

- `docs/BUNDLE_LAYOUT_v1.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`

---

## 10. Replay v1

Replay verification reconstructs the deterministic state and verifies it against the expected hash.

Replay must be deterministic.

Replay must not depend on:

- current time
- random values
- network responses
- machine-specific ordering
- non-pinned external dependencies
- mutable environment data

The replay and verification procedure is defined in:

- `docs/VERIFY_SPEC_v1.0.md`

---

## 11. Verification Result

A DigiEmu verification result must produce a clear machine-readable result.

The minimal semantic outcomes are:

```text
PASS
FAIL
ERROR
```

Where:

- `PASS` means the reconstructed hash matches the expected hash.
- `FAIL` means verification completed but the reconstructed hash does not match.
- `ERROR` means verification could not be completed due to missing input, invalid format or runtime failure.

The result schema is defined in:

- `docs/VERIFY_RESULT_SCHEMA_v1.json`

---

## 12. Conformance Levels

DigiEmu Core defines three initial conformance levels.

### DigiEmu Core Reader

A Reader can:

- read a DigiEmu bundle
- parse snapshot files
- detect expected hash fields
- inspect verification reports

### DigiEmu Core Verifier

A Verifier can:

- reconstruct the deterministic state
- produce canonical JSON
- compute SHA-256 hashes
- compare computed and expected hashes
- return PASS / FAIL / ERROR results

### DigiEmu Core Producer

A Producer can:

- generate valid deterministic snapshots
- create valid expected hash files
- separate inside-hash and outside-hash data correctly
- produce portable verification bundles

A full conforming implementation should satisfy Reader, Verifier and Producer requirements.

---

## 13. Security Considerations

DigiEmu Core improves auditability and reproducibility, but it does not automatically secure the full system.

Important risks include:

- tampered snapshots
- altered expected hashes
- manipulated verification reports
- ambiguous state boundaries
- non-deterministic replay dependencies
- false confidence caused by incomplete state capture
- confusion between reproducibility and correctness

Security considerations are further documented in:

- `docs/THREAT_MODEL.md`
- `SECURITY.md`

Future DigiEmu Secure work may define stronger protections for signatures, key management, tamper-evident storage and trusted verifier operation.

---

## 14. Governance and Versioning

DigiEmu specifications should be versioned independently from implementation releases.

This document is a public review draft.

Recommended interpretation:

```text
Implementation baseline: v1.0.0
Public standard draft: DigiEmu Core Specification v0.9
Future milestone: DigiEmu Core Specification v1.0
```

Versioning policy is defined in:

- `docs/VERSIONING_POLICY_v1.0.md`

---

## 15. Relationship to DigiEmu Proof

DigiEmu Core defines deterministic state identity.

DigiEmu Proof verifies transitions and chains between states.

In simple terms:

- Core defines the boundary.
- Proof verifies the transition.
- Secure protects the evidence.
- Enterprise operationalizes the stack.

---

## 16. Reference Documents

This specification consolidates the following repository documents:

- `docs/SPEC_INDEX_v1.0.md`
- `docs/SNAPSHOT_HASH_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
- `docs/BUNDLE_LAYOUT_v1.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/VERIFY_RESULT_SCHEMA_v1.json`
- `docs/VERSIONING_POLICY_v1.0.md`
- `docs/PUBLIC_API_v1.0.md`
- `docs/CLI_CONTRACT_v1.0.md`
- `docs/CLI_VERIFY_v1.0.md`
- `docs/THREAT_MODEL.md`
- `SECURITY.md`

---

## 17. Test Vectors

A future v1.0 specification should include normative test vectors.

Each test vector should define:

- input snapshot
- canonical JSON output
- expected SHA-256 hash
- expected verification result
- tampered variant
- expected FAIL result

Test vectors are required before DigiEmu Core can be considered a stable external standard.

---

## 18. Current Status

DigiEmu Core Specification v0.9 is a public review draft.

It is intended to consolidate the current implementation baseline and prepare the path toward a future DigiEmu Core Specification v1.0.

The next milestone is external review, test vectors and conformance validation.

