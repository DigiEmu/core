# Core 2.0 Conformance Testdata

This folder contains minimal test cases for Core 2.0 conformance testing. Each case contains an `input.json` and an `expected_verify_result.json` representing the expected high-level outcome.

Structure
---------
- `basic_pass/` — Minimal input leading to `PASS`.
- `hash_mismatch_fail/` — Input leading to `FAIL` due to a mismatched hash.
- `invalid_schema_error/` — Input that should trigger an `ERROR` due to invalid schema.

Usage
-----
Copy the directory into your test harness or run the checks locally against your implementation. See `docs/CORE_2_CONFORMANCE_PACK.md` for guidance.
