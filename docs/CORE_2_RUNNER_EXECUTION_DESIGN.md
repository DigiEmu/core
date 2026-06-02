# Core 2.0 Runner Real Execution Design

Status: design / draft / Phase 2 input parsing started

Milestone: post-v2.0.0-draft.2

Purpose: Define future runner behavior before implementation.

This document describes how the Core 2.0 conformance runner can evolve from structural validation of expected verify results toward real input execution and expected-result comparison. This is a design step only and does not change production code, CLI behavior, schemas, testdata, tags, or v1.0 compatibility.

Current state
-------------

- The current runner discovers conformance cases.
- The current runner validates `expected_verify_result.json` structurally.
- The current runner reports `case_passed` when the expected verify result is structurally valid.
- Some current cases are representative partner test vectors, not fully executed verification scenarios yet.

Target state
------------

- Runner reads `input.json`.
- Runner executes or simulates the relevant verify path.
- Runner produces an actual verify result.
- Runner compares actual result with `expected_verify_result.json`.
- Runner reports both case structural validity and execution match status.
- JSON report remains machine-readable and CI-friendly.

Execution phases
----------------

### Phase 1: Structural validation

Status: already implemented

Behavior:
- Validate `expected_verify_result.json`.
- Report `case_passed` for structurally valid expected results.

### Phase 2: Input parsing

Status: implemented for JSON parsing and malformed JSON classification

Behavior:
- Read `input.json`.
- Fail malformed JSON cases deterministically.
- Map malformed JSON to `ERROR` / `INVALID_SNAPSHOT_SCHEMA`.
- Preserve representative expected-result behavior for well-formed inputs until full verify execution exists.

### Phase 3: Verify execution

Status: future

Behavior:
- Run canonicalization/reconstruction/hash validation where representable.
- Produce actual Verify Result v2.
- Compare actual result to expected result.

### Phase 4: Partner implementation comparison

Status: future

Behavior:
- Allow partner-produced verify results to be checked against expected outputs.
- Support conformance automation for external implementations.

Report model notes
------------------

- Keep `report_version` stable for the current draft unless schema changes are required.
- Consider adding `actual_result` and expected-result comparison fields in a future report schema.
- Do not change the current JSON report schema in this design step.

Reason-code implications
------------------------

- Real execution will clarify `INTERNAL_ERROR` boundaries.
- Phase 2 currently maps malformed JSON to `INVALID_SNAPSHOT_SCHEMA`; future review may decide whether a separate `INVALID_JSON` code is needed.
- Real execution will enable stronger tests for `UNSUPPORTED_PROFILE` and `MISSING_REFERENCE`.

Non-goals
---------

- No stable CLI promotion.
- No HTTP API server implementation.
- No signing or key custody.
- No policy judgement.
- No v1.0 behavior change.
- No authentication or tenancy implementation.
- No tag creation, movement, or promotion.

Recommended next PRs
--------------------

1. Extend Phase 2 beyond JSON syntax into deterministic input shape/schema classification.
2. Add actual-vs-expected comparison design for future report schema.
3. Add representable `UNSUPPORTED_PROFILE` and `MISSING_REFERENCE` cases.
4. Scope Phase 3 verify execution without changing v1.0 behavior.

Compatibility statement
-----------------------

This document is a design draft for future Core 2.0 runner behavior. It does not promote Core 2.0 draft features to stable, does not change v1.0 behavior, and does not implement any runner, CLI, schema, HTTP API, signing, authentication, tenancy, policy, or key-custody changes.

Change history
--------------

- 2026-06-02 — Added runner real execution design for post-`v2.0.0-draft.2` stabilization planning.
- 2026-06-02 — Marked Phase 2 JSON parsing and malformed JSON classification as started.
