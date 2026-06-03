# Core 2.0 Actual-vs-Expected Comparison Design

Status
------

- Milestone: post-v2.0.0-draft.3
- Stability: design / draft / Phase 3b runner integration started
- Purpose: Define future runner comparison semantics before implementation.

Purpose
-------

This document defines future Core 2.0 conformance runner actual-vs-expected comparison semantics before implementation. It is documentation-only and does not change production code, schemas, testdata, CLI behavior, tags, or v1.0 compatibility.

Current state
-------------

- The runner discovers conformance cases.
- The runner reads and parses `input.json`.
- Malformed input JSON is handled deterministically.
- Missing conformance case files are covered by tests.
- The official conformance pack reports total=10 passed=10 failed=0.
- Full verify execution is not implemented yet.

Target state
------------

- The runner produces an observed verify result.
- The runner loads `expected_verify_result.json`.
- The runner compares observed result with expected result.
- The runner marks a case as passed only when observed and expected semantics match.
- The runner preserves machine-readable JSON report output.
- The runner remains under the experimental namespace until explicitly stabilized.

Comparison fields
-----------------

Required for initial MVP:

- `result`
- `reason_code`

Optional future fields:

- `profile`
- `hash_algorithm`
- `hash`
- `error_details`
- `schema_validation_details`

Comparison rules
----------------

- `PASS` must match `PASS`.
- `FAIL` must match `FAIL`.
- `ERROR` must match `ERROR`.
- `reason_code` must match exactly for initial MVP.
- Unknown or unsupported expected reason codes remain schema validation failures.
- Malformed input observed result should continue to map to `ERROR` / `INVALID_SNAPSHOT_SCHEMA`.
- Missing `expected_verify_result.json` remains a structural conformance-case failure.

Phased implementation
---------------------

### Phase 3a: Internal observed-result model

Status: implemented for internal `result` / `reason_code` comparison helper

Behavior:

- Introduce internal observed verify result representation.
- Do not change public report shape.
- Compare observed result and `reason_code` to expected result.

### Phase 3b: Actual-vs-expected unit tests

Status: implemented for helper-level and malformed-input runner integration coverage

Behavior:

- Add tests for matched and mismatched result/reason_code.
- Keep official 10-case pack passing.

### Phase 3c: Report evolution consideration

Behavior:

- Evaluate whether future reports should expose `expected_result` and `observed_result` separately.
- Do not change report schema until needed.

### Phase 4: Full verify execution

Behavior:

- Replace representative observed results with real verification output where implemented.
- Keep backward compatibility rules explicit.

Reason-code implications
------------------------

- Exact `reason_code` comparison will make reason code semantics stricter.
- `INTERNAL_ERROR` boundaries may need refinement.
- `INVALID_SNAPSHOT_SCHEMA` may later split into `INVALID_JSON` if partner feedback requires it.
- `UNSUPPORTED_PROFILE` and `MISSING_REFERENCE` should become stronger candidates for representable execution cases.

Non-goals
---------

- No stable CLI promotion.
- No HTTP API server.
- No signing or key custody.
- No policy judgement.
- No v1.0 behavior change.
- No report schema change in this design step.

Recommended next PRs
--------------------

- Expand observed-result comparison coverage where new observed results are introduced.
- Keep public JSON report shape unchanged.
- Add representable `UNSUPPORTED_PROFILE` and `MISSING_REFERENCE` cases after comparison helper exists.

Change history
--------------

- 2026-06-03 — Added actual-vs-expected comparison design for post-draft.3 runner hardening.
- 2026-06-03 — Marked Phase 3a internal `result` / `reason_code` comparison helper as implemented.
- 2026-06-03 — Marked Phase 3b runner integration coverage as started for observed-result paths.

Compatibility statement
-----------------------

This design does not implement full verify execution or promote Core 2.0 to stable. Experimental runner behavior remains draft-facing, report schema changes are deferred until explicitly needed, and v1.0 behavior remains unchanged.
