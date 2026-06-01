# DigiEmu Core 2.0 — Conformance Pack (Draft)

Status: DRAFT

Purpose
-------
This document defines a Core 2.0 Conformance Pack for external implementations. The Conformance Pack provides a small, normative set of test cases and rules that allow implementers to verify that their implementation produces deterministic canonical snapshots, consistent hashes, and valid `Verify Result` v2 outputs without changing v1.0 behavior.

Scope and constraints
---------------------
- Conformance tests and artifacts are informational and opt-in. They do not change or replace any v1.0 contracts.  
- Conformance MUST be deterministic and MUST NOT require user management, signing, key custody, orchestration, or policy judgement.
- Conformance focuses on canonicalization, snapshot formation, hash generation, reconstruction, and verify-result semantics.

Conformance model
-----------------
1. Determinism: For the same input, a conforming implementation MUST produce the same canonical snapshot.
2. Hash equality: For the same canonical snapshot, a conforming implementation MUST produce the same canonical snapshot hash (given the same canonicalization profile and hash algorithm).
3. Verify result structure: Implementations MUST be able to produce a `Verify Result` v2 structured output matching `schemas/verify_result_v2.schema.json` (draft).

Test cases
----------
The conformance testdata includes three minimal cases:
- `basic_pass`: minimal deterministic input that SHOULD produce `PASS` + `STATE_RECONSTRUCTED`.
- `hash_mismatch_fail`: input crafted so that the computed snapshot hash differs from the expected value; expected `FAIL` + `HASH_MISMATCH`.
- `invalid_schema_error`: input that violates the snapshot input schema, expected `ERROR` + `INVALID_SNAPSHOT_SCHEMA`.

How to use
----------
1. Copy the `testdata/core_2_conformance` directory to your implementation test harness.  
2. Run your canonicalization and hashing logic on the `input.json` for each case.  
3. Compare produced `snapshot_hash` with `expected_verify_result.json` (or validate produced `Verify Result` v2 against `schemas/verify_result_v2.schema.json`).

Reporting
---------
Implementers SHOULD provide a short report summarizing: input used, produced canonical snapshot, computed snapshot hash, the produced `Verify Result` v2 document (if any), and a pass/fail status per case.

Acceptance
----------
Conformance is considered satisfied when the implementation reproduces the deterministic canonical snapshot and hash for `basic_pass` and emits the expected `Verify Result` v2 structure for all three cases.
