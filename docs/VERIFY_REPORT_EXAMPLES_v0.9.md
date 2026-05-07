# DigiEmu Core Verify Report Examples v0.9

Status: Public Review Draft  
Specification Layer: DigiEmu Core  
Version: v0.9  
Implementation Baseline: v1.0.0  

---

## 1. Purpose

This document provides machine-readable examples of DigiEmu Core verification reports.

A verification report records the result of comparing a deterministic canonical snapshot hash with an expected hash.

The goal of this document is to make DigiEmu Core verification output:

- reproducible
- machine-readable
- implementation-friendly
- audit-ready
- suitable for conformance testing

This document complements:

- DigiEmu Core Specification v0.9
- DigiEmu Core Test Vectors v0.9
- DigiEmu Core Conformance v0.9
- Snapshot Hash v1.0
- Verify Spec v1.0

---

## 2. Verification Report Result Types

A DigiEmu Core verification report SHOULD use one of the following result values:

```text
PASS
FAIL
ERROR
```

### PASS

`PASS` means that the computed hash of the canonical snapshot is identical to the expected hash.

### FAIL

`FAIL` means that the snapshot was parsed and canonicalized successfully, but the computed hash does not match the expected hash.

### ERROR

`ERROR` means that verification could not be completed because the input was invalid, incomplete, malformed, unsupported, or otherwise not processable.

---

## 3. Common Report Fields

A verification report SHOULD include the following fields where applicable:

```json
{
  "result": "PASS",
  "computed_hash": "...",
  "expected_hash": "...",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256"
}
```

Recommended fields:

| Field | Description |
|---|---|
| `result` | Verification result: `PASS`, `FAIL`, or `ERROR` |
| `computed_hash` | Hash computed from the canonical snapshot |
| `expected_hash` | Hash expected by the bundle, test vector, or verification target |
| `canonicalization` | Canonicalization method used |
| `hash_algorithm` | Hash algorithm used |
| `reason` | Machine-readable reason for `FAIL` or `ERROR` |
| `case_id` | Optional identifier of the verified case |
| `test_vector_id` | Optional test vector identifier |
| `spec_version` | Optional DigiEmu Core specification version |
| `verified_at` | Optional timestamp outside the deterministic hash boundary |

---

## 4. PASS Example

This example uses Test Vector 001 from DigiEmu Core Test Vectors v0.9.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "low"
}
```

### Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v1","risk_level":"low"}
```

### Expected SHA-256

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### Verification Report

```json
{
  "result": "PASS",
  "case_id": "CORE-TV-001",
  "test_vector_id": "CORE-TV-001",
  "computed_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256",
  "spec_version": "DigiEmu Core Specification v0.9"
}
```

### Interpretation

The canonical snapshot was reconstructed successfully.

The computed hash matches the expected hash.

The verification result is therefore:

```text
PASS
```

---

## 5. FAIL Example

A `FAIL` report is produced when the snapshot is valid and can be canonicalized, but the computed hash differs from the expected hash.

This usually means that the deterministic input state has changed.

Examples:

- changed decision value
- changed risk level
- changed policy reference
- changed symptom list
- changed vital value
- changed deterministic field inside the hash boundary

### Tampered Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "deny",
  "policy": "policy_v1",
  "risk_level": "low"
}
```

### Tampered Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"deny","policy":"policy_v1","risk_level":"low"}
```

### Computed SHA-256

```text
547cfdfc0f1cbd4074b0df84a70a2dcf9f1ca9223535a0c46b8bda4372d9f0f0
```

### Expected SHA-256

The expected hash remains the original hash from Test Vector 001:

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### Verification Report

```json
{
  "result": "FAIL",
  "case_id": "CORE-TV-001",
  "test_vector_id": "CORE-TV-001",
  "computed_hash": "547cfdfc0f1cbd4074b0df84a70a2dcf9f1ca9223535a0c46b8bda4372d9f0f0",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256",
  "reason": "computed_hash_does_not_match_expected_hash",
  "spec_version": "DigiEmu Core Specification v0.9"
}
```

### Interpretation

The snapshot was valid and could be canonicalized.

However, the computed hash does not match the expected hash.

The verification result is therefore:

```text
FAIL
```

---

## 6. ERROR Example

An `ERROR` report is produced when verification cannot be completed.

This is different from `FAIL`.

A `FAIL` result means:

```text
The input was processable, but the hash did not match.
```

An `ERROR` result means:

```text
The input could not be processed correctly.
```

Examples:

- malformed JSON
- missing snapshot
- missing expected hash
- unsupported hash algorithm
- unsupported canonicalization method
- invalid bundle structure
- invalid field type
- unreadable artifact

### Invalid Input Example

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "low"
```

This JSON is malformed because the closing brace is missing.

### Verification Report

```json
{
  "result": "ERROR",
  "reason": "malformed_json",
  "message": "Verification could not be completed because the input snapshot is not valid JSON.",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256",
  "spec_version": "DigiEmu Core Specification v0.9"
}
```

### Interpretation

The verifier could not parse the input snapshot.

No deterministic canonical snapshot could be reconstructed.

The verification result is therefore:

```text
ERROR
```

---

## 7. ERROR Example: Missing Expected Hash

A verifier may also return `ERROR` if the snapshot is valid but the expected hash is missing.

### Verification Report

```json
{
  "result": "ERROR",
  "case_id": "CORE-TV-001",
  "reason": "missing_expected_hash",
  "message": "Verification could not be completed because no expected hash was provided.",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256",
  "spec_version": "DigiEmu Core Specification v0.9"
}
```

---

## 8. Machine-Readable Reason Codes

Implementations SHOULD use stable machine-readable reason codes.

Recommended reason codes:

```text
computed_hash_does_not_match_expected_hash
malformed_json
missing_snapshot
missing_expected_hash
unsupported_hash_algorithm
unsupported_canonicalization
invalid_bundle_layout
invalid_field_type
unreadable_artifact
internal_verifier_error
```

Reason codes SHOULD be lowercase.

Reason codes SHOULD use snake_case.

Reason codes SHOULD remain stable across implementations where possible.

---

## 9. Deterministic Boundary Note

Verification reports themselves are usually outside the deterministic snapshot hash boundary.

The report may contain non-deterministic metadata such as:

- verification timestamp
- verifier implementation name
- verifier version
- runtime environment
- human-readable messages
- audit notes

These fields SHOULD NOT change the hash of the verified snapshot.

The deterministic hash is computed from the canonical snapshot, not from the verification report.

---

## 10. Optional Metadata Example

Implementations MAY include additional metadata outside the deterministic hash boundary.

Example:

```json
{
  "result": "PASS",
  "case_id": "CORE-TV-001",
  "test_vector_id": "CORE-TV-001",
  "computed_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256",
  "spec_version": "DigiEmu Core Specification v0.9",
  "metadata": {
    "verified_at": "2026-05-07T00:00:00Z",
    "verifier": "digiemu-core-reference-verifier",
    "verifier_version": "v1.0.0"
  }
}
```

The `metadata` object is audit context.

It SHOULD NOT be used as part of the deterministic snapshot hash.

---

## 11. Minimal Valid Reports

### Minimal PASS Report

```json
{
  "result": "PASS",
  "computed_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023"
}
```

### Minimal FAIL Report

```json
{
  "result": "FAIL",
  "computed_hash": "547cfdfc0f1cbd4074b0df84a70a2dcf9f1ca9223535a0c46b8bda4372d9f0f0",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "reason": "computed_hash_does_not_match_expected_hash"
}
```

### Minimal ERROR Report

```json
{
  "result": "ERROR",
  "reason": "malformed_json"
}
```

---

## 12. Conformance Relevance

A DigiEmu Core Verifier SHOULD be able to produce verification reports equivalent to the examples in this document.

A conforming verifier SHOULD distinguish clearly between:

```text
PASS  = computed hash matches expected hash
FAIL  = computed hash does not match expected hash
ERROR = verification could not be completed
```

A verifier MUST NOT return `PASS` if the computed hash differs from the expected hash.

A verifier SHOULD NOT return `FAIL` for malformed or unprocessable input.

Malformed or unprocessable input SHOULD result in `ERROR`.

---

## 13. Current Status

This document is a Public Review Draft.

It is intended to support early implementers, reviewers, auditors, and conformance tooling.

Future versions may define:

- a formal JSON Schema for verification reports
- stricter required fields
- versioned reason-code registries
- compatibility rules for verifier implementations
- machine-readable conformance declarations