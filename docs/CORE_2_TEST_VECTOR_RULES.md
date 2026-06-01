# Core 2.0 Test Vector Rules (Draft)

Status: DRAFT

Purpose
-------
Define the rules and expectations for Core 2.0 test vectors used in conformance testing.

Principles
----------
- Determinism: Test vectors MUST be deterministic and self-contained.
- Non-destructive: Test vectors MUST NOT require modifications to production schemas or code.
- Versioned: New vectors for v2 MUST be stored separately from v1 vectors.

File layout
-----------
- Each case is a directory under `testdata/core_2_conformance/<case_name>/` containing:
  - `input.json` — canonical input to the Core canonicalizer/reconstructor.
  - `expected_verify_result.json` — expected `Verify Result` v2 document for the case.

Rules
-----
- `input.json` MUST be minimal and focused on the case being tested.
- `expected_verify_result.json` MUST conform to `schemas/verify_result_v2.schema.json` (draft).
- PASS examples SHOULD include `snapshot_hash` representing the expected canonical snapshot hash.
- FAIL examples SHOULD include both `snapshot_hash_expected` and `snapshot_hash_actual` where applicable.
- ERROR examples MAY omit `snapshot_hash` if a valid snapshot cannot be produced.
- `failure_path` SHOULD use JSON Pointer syntax (RFC 6901) when pointing to failing data elements.

Maintenance
-----------
- Adding or modifying test vectors MUST be recorded in the project changelog and linked from the conformance pack.
