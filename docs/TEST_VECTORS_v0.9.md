# DigiEmu Core Test Vectors v0.9

Status: Public Review Draft  
Version: 0.9  
Repository: DigiEmu/core  
Scope: Minimal reproducible test vectors for DigiEmu Core deterministic snapshot verification.

---

## 1. Purpose

This document defines minimal test vectors for DigiEmu Core Specification v0.9.

The goal is to provide simple, reproducible examples that allow implementers to verify whether their canonicalization, hashing and verification behavior is compatible with DigiEmu Core.

A valid DigiEmu Core test vector should demonstrate:

- deterministic input snapshot
- canonical JSON representation
- expected SHA-256 hash
- expected verification result
- tampered variant
- expected verification failure

---

## 2. Test Vector Model

Each test vector contains:

1. input snapshot
2. canonical JSON
3. expected SHA-256 hash
4. expected PASS result
5. tampered snapshot
6. expected FAIL result

The test vectors are intentionally small. They are designed for clarity, not for domain completeness.

---

## 3. Canonicalization Assumptions

For these v0.9 test vectors, canonical JSON is represented as:

- UTF-8 encoded JSON
- no insignificant whitespace
- stable object key ordering
- deterministic array order
- exact string values
- deterministic number representation

Example canonical form:

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v1","risk_level":"low"}
```

The SHA-256 hash is computed over the exact canonical JSON byte representation.

---

## 4. Test Vector 001 — Minimal Decision State

### 4.1 Input Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "policy": "policy_v1",
  "decision": "allow",
  "risk_level": "low"
}
```

### 4.2 Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"allow","policy":"policy_v1","risk_level":"low"}
```

### 4.3 Expected SHA-256

```text
bcc9c8ee5e6e269d598adaa26cb4e3875726da29b3568ab83f05001bae2ce023
```

### 4.4 Expected Verification Result

```json
{
  "result": "PASS",
  "reason": "computed_hash_matches_expected_hash"
}
```

### 4.5 Tampered Snapshot

```json
{
  "case_id": "CORE-TV-001",
  "policy": "policy_v1",
  "decision": "deny",
  "risk_level": "low"
}
```

### 4.6 Tampered Canonical JSON

```json
{"case_id":"CORE-TV-001","decision":"deny","policy":"policy_v1","risk_level":"low"}
```

### 4.7 Expected Tampered Verification Result

```json
{
  "result": "FAIL",
  "reason": "computed_hash_does_not_match_expected_hash"
}
```

---

## 5. Test Vector 002 — Medical Triage Style State

### 5.1 Input Snapshot

```json
{
  "case_id": "GZC-CASE-001",
  "symptoms": [
    "cough",
    "fever"
  ],
  "red_flags": {
    "chest_pain": false,
    "confusion": false,
    "severe_breathing_difficulty": false,
    "unconscious": false
  },
  "vitals": {
    "oxygen_saturation_percent": 97,
    "respiratory_rate_per_min": 18
  },
  "triage_level": "low"
}
```

### 5.2 Canonical JSON

```json
{"case_id":"GZC-CASE-001","red_flags":{"chest_pain":false,"confusion":false,"severe_breathing_difficulty":false,"unconscious":false},"symptoms":["cough","fever"],"triage_level":"low","vitals":{"oxygen_saturation_percent":97,"respiratory_rate_per_min":18}}
```

### 5.3 Expected SHA-256

```text
e796de537dab5cd5a42f2258fc2c0536ea9f210a32910cf416708ad54e3210f8
```

### 5.4 Expected Verification Result

```json
{
  "result": "PASS",
  "reason": "computed_hash_matches_expected_hash"
}
```

### 5.5 Tampered Snapshot

```json
{
  "case_id": "GZC-CASE-001",
  "symptoms": [
    "cough",
    "fever"
  ],
  "red_flags": {
    "chest_pain": false,
    "confusion": false,
    "severe_breathing_difficulty": true,
    "unconscious": false
  },
  "vitals": {
    "oxygen_saturation_percent": 97,
    "respiratory_rate_per_min": 18
  },
  "triage_level": "low"
}
```

### 5.6 Tampered Canonical JSON

```json
{"case_id":"GZC-CASE-001","red_flags":{"chest_pain":false,"confusion":false,"severe_breathing_difficulty":true,"unconscious":false},"symptoms":["cough","fever"],"triage_level":"low","vitals":{"oxygen_saturation_percent":97,"respiratory_rate_per_min":18}}
```

### 5.7 Expected Tampered Verification Result

```json
{
  "result": "FAIL",
  "reason": "computed_hash_does_not_match_expected_hash"
}
```

---

## 6. Notes

These test vectors are part of the DigiEmu Core v0.9 public review draft.

The expected hashes are placeholders until they are computed and confirmed by the reference implementation.

Before DigiEmu Core Specification v1.0, this document should be updated with:

- final canonicalization rules
- computed SHA-256 values
- reference verifier output
- CLI commands for reproducing each result

---

## 7. Relationship to Core Specification

This document supports:

- `docs/DIGIEMU_CORE_SPEC_v0.9.md`
- `docs/SNAPSHOT_HASH_v1.0.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/VERIFY_RESULT_SCHEMA_v1.json`

The test vectors provide concrete examples for the conformance and verification concepts defined in the Core Specification.