# DigiEmu Core Negative Test Vectors v0.9

Status: Public Review Draft  
Specification Layer: DigiEmu Core  
Version: v0.9  
Implementation Baseline: v1.0.0  

---

## 1. Purpose

This document defines negative test vectors for DigiEmu Core verification.

Negative test vectors are intentionally invalid, tampered, incomplete, malformed, or unsupported inputs.

The goal is to ensure that DigiEmu Core implementations correctly distinguish between:

```text
PASS  = computed hash matches expected hash
FAIL  = input is processable, but computed hash does not match expected hash
ERROR = verification could not be completed
```

This document complements:

- DigiEmu Core Specification v0.9
- DigiEmu Core Test Vectors v0.9
- DigiEmu Core Conformance v0.9
- DigiEmu Core Conformance Declaration v0.9
- DigiEmu Core Verify Report Examples v0.9
- DigiEmu Core Verify Report Schema v0.9

---

## 2. Negative Verification Outcomes

DigiEmu Core negative test vectors are divided into two categories:

```text
FAIL test vectors
ERROR test vectors
```

### FAIL Test Vectors

A `FAIL` test vector is valid, parseable, and canonicalizable.

The verifier can compute a hash, but the computed hash does not match the expected hash.

Typical causes:

- deterministic snapshot field changed
- decision changed
- risk level changed
- policy reference changed
- symptom list changed
- vital value changed
- deterministic inside-hash content changed

### ERROR Test Vectors

An `ERROR` test vector cannot be verified because verification cannot be completed.

Typical causes:

- malformed JSON
- missing snapshot
- missing expected hash
- unsupported hash algorithm
- unsupported canonicalization method
- invalid bundle layout
- invalid field type
- unreadable artifact

---

## 3. Base Positive Test Vector

The negative test vectors in this document use Test Vector 001 as a baseline.

### Base Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "low"
}
```

### Base Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v1","risk_level":"low"}
```

### Base Expected SHA-256

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

The base vector is expected to produce:

```text
PASS
```

The following negative vectors intentionally deviate from this baseline.

---

## 4. FAIL Test Vector 001 — Changed Decision

Identifier:

```text
CORE-NTV-FAIL-001
```

Purpose:

Verify that changing a deterministic decision field causes verification to fail.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "deny",
  "policy": "policy_v1",
  "risk_level": "low"
}
```

### Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"deny","policy":"policy_v1","risk_level":"low"}
```

### Expected SHA-256

The expected hash remains the original hash from Test Vector 001:

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### Computed SHA-256

```text
547cfdfc0f1cbd4074b0df84a70a2dcf9f1ca9223535a0c46b8bda4372d9f0f0
```

### Expected Result

```text
FAIL
```

### Expected Reason Code

```text
computed_hash_does_not_match_expected_hash
```

---

## 5. FAIL Test Vector 002 — Changed Risk Level

Identifier:

```text
CORE-NTV-FAIL-002
```

Purpose:

Verify that changing a deterministic risk field causes verification to fail.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "high"
}
```

### Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v1","risk_level":"high"}
```

### Expected SHA-256

The expected hash remains the original hash from Test Vector 001:

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### Expected Result

```text
FAIL
```

### Expected Reason Code

```text
computed_hash_does_not_match_expected_hash
```

Note:

The exact computed hash MAY be included by implementations in their verification report.

---

## 6. FAIL Test Vector 003 — Changed Policy Reference

Identifier:

```text
CORE-NTV-FAIL-003
```

Purpose:

Verify that changing a deterministic policy reference causes verification to fail.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v2",
  "risk_level": "low"
}
```

### Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v2","risk_level":"low"}
```

### Expected SHA-256

The expected hash remains the original hash from Test Vector 001:

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### Expected Result

```text
FAIL
```

### Expected Reason Code

```text
computed_hash_does_not_match_expected_hash
```

---

## 7. ERROR Test Vector 001 — Malformed JSON

Identifier:

```text
CORE-NTV-ERROR-001
```

Purpose:

Verify that malformed JSON results in `ERROR`, not `FAIL`.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "low"
```

This JSON is malformed because the closing brace is missing.

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
malformed_json
```

---

## 8. ERROR Test Vector 002 — Missing Snapshot

Identifier:

```text
CORE-NTV-ERROR-002
```

Purpose:

Verify that a missing snapshot results in `ERROR`.

### Input

```text
No snapshot artifact is provided.
```

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
missing_snapshot
```

---

## 9. ERROR Test Vector 003 — Missing Expected Hash

Identifier:

```text
CORE-NTV-ERROR-003
```

Purpose:

Verify that a valid snapshot without an expected hash results in `ERROR`.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": "low"
}
```

### Expected Hash

```text
missing
```

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
missing_expected_hash
```

---

## 10. ERROR Test Vector 004 — Unsupported Hash Algorithm

Identifier:

```text
CORE-NTV-ERROR-004
```

Purpose:

Verify that an unsupported hash algorithm results in `ERROR`.

### Input Metadata

```json
{
  "hash_algorithm": "MD5"
}
```

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
unsupported_hash_algorithm
```

---

## 11. ERROR Test Vector 005 — Unsupported Canonicalization Method

Identifier:

```text
CORE-NTV-ERROR-005
```

Purpose:

Verify that an unsupported canonicalization method results in `ERROR`.

### Input Metadata

```json
{
  "canonicalization": "custom_json_v0"
}
```

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
unsupported_canonicalization
```

---

## 12. ERROR Test Vector 006 — Invalid Field Type

Identifier:

```text
CORE-NTV-ERROR-006
```

Purpose:

Verify that invalid field types may result in `ERROR` if they violate the expected snapshot structure.

### Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "decision": "allow",
  "policy": "policy_v1",
  "risk_level": 1
}
```

### Expected Result

```text
ERROR
```

### Expected Reason Code

```text
invalid_field_type
```

Note:

If an implementation treats this as a processable generic JSON snapshot, it may compute a hash and return `FAIL`.

However, implementations that validate the expected Test Vector 001 structure SHOULD return `ERROR`.

---

## 13. Expected Reason Codes

Negative test vectors SHOULD use stable machine-readable reason codes.

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

---

## 14. Conformance Relevance

A DigiEmu Core Verifier SHOULD pass both positive and negative test vectors.

A verifier claiming DigiEmu Core Verifier support SHOULD be able to distinguish:

```text
valid unchanged snapshot     → PASS
valid tampered snapshot      → FAIL
malformed or incomplete input → ERROR
```

A conforming verifier MUST NOT return `PASS` for a tampered snapshot.

A conforming verifier SHOULD NOT return `FAIL` for malformed JSON.

A conforming verifier SHOULD return `ERROR` when verification cannot be completed.

---

## 15. Relationship to Test Vectors v0.9

This document extends:

```text
docs/TEST_VECTORS_v0.9.md
```

The main Test Vectors document defines positive reproducibility examples.

This document defines negative behavior expectations.

Together they support:

```text
positive reproducibility testing
negative failure testing
error handling testing
```

---

## 16. Relationship to Verify Report Schema v0.9

Negative test vectors can be represented using DigiEmu Core verification reports.

For example, a FAIL vector should produce a report similar to:

```json
{
  "result": "FAIL",
  "computed_hash": "547cfdfc0f1cbd4074b0df84a70a2dcf9f1ca9223535a0c46b8bda4372d9f0f0",
  "expected_hash": "bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023",
  "reason": "computed_hash_does_not_match_expected_hash",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256"
}
```

An ERROR vector should produce a report similar to:

```json
{
  "result": "ERROR",
  "reason": "malformed_json",
  "canonicalization": "canonical_json_v1",
  "hash_algorithm": "SHA-256"
}
```

These reports SHOULD validate against:

```text
docs/VERIFY_REPORT_SCHEMA_v0.9.json
```

---

## 17. Current Status

This document is a Public Review Draft.

Future versions may define:

- exact computed hashes for all FAIL vectors
- machine-readable negative test vector files
- JSON Schema for test vector manifests
- implementation test harness expectations
- stricter ERROR behavior requirements
- compatibility with DigiEmu Secure profiles