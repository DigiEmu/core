# Core 2.0 Reason Code Review and Stabilization Notes

Status: draft stabilization review

This document records the Core 2.0 reason-code review for the current draft milestone. It is documentation-only and does not change production code, schemas, testdata, CLI behavior, tags, or compatibility guarantees.

Review scope
------------

- `docs/CORE_2_REASON_CODES.md`
- `docs/CORE_2_VERIFY_RESULT_v2_DRAFT.md`
- `schemas/verify_result_v2.schema.json`
- `testdata/core_2_conformance/`
- Experimental conformance runner report fields that surface `reason_code`
- Future runner real execution design in `docs/CORE_2_RUNNER_EXECUTION_DESIGN.md`

Review outcome
--------------

- The current reason-code set is suitable for continued Core 2.0 draft partner review.
- No immediate schema or testdata changes are required for this stabilization note.
- The registry, Verify Result v2 draft, schema enum, and conformance fixtures appear aligned around the same draft code set.
- Core 2.0 is not promoted to stable by this review.

Current reason-code set
-----------------------

PASS:
- `STATE_RECONSTRUCTED`

FAIL:
- `HASH_MISMATCH`
- `RECONSTRUCTION_RULE_MISMATCH`
- `CHAIN_BROKEN`
- `TRANSITION_INVALID`
- `INSIDE_HASH_FIELD_CHANGED`

ERROR:
- `CANONICALIZATION_ERROR`
- `INVALID_SNAPSHOT_SCHEMA`
- `MISSING_REQUIRED_FIELD`
- `UNSUPPORTED_HASH_ALGORITHM`
- `UNSUPPORTED_CANONICALIZATION_PROFILE`
- `UNSUPPORTED_RECONSTRUCTION_PROFILE`
- `INTERNAL_ERROR`

Stabilization notes
-------------------

- Reason codes are intended to be stable, machine-readable identifiers for `Verify Result` v2 consumers.
- Codes are grouped by `result` to keep PASS, FAIL, and ERROR handling deterministic.
- Existing codes should not be renamed once partner integrations depend on them.
- New codes should be added only when an existing code would create ambiguous consumer behavior.
- Human-readable wording may evolve, but code identifiers should remain stable after a stable Core 2.0 release.

Release risks
-------------

- `INTERNAL_ERROR` is intentionally broad and should remain a last-resort fallback, not a substitute for specific validation errors.
- The distinction between `FAIL` and `ERROR` must remain clear for partners: `FAIL` means verification completed with a negative result; `ERROR` means verification could not complete normally.
- Adding codes without updating the registry, Verify Result v2 draft, schema enum, examples, and conformance cases would create compatibility drift.
- Partner automation may begin depending on draft identifiers before stable release, so draft/stable boundaries must remain explicit.
- Real runner execution will help validate whether reason-code semantics are specific enough for actual-vs-expected comparison.

Execution implications
----------------------

- Real execution will clarify `INTERNAL_ERROR` boundaries by separating implementation failures from expected verification errors.
- Current runner input parsing maps malformed JSON to `INVALID_SNAPSHOT_SCHEMA`.
- Future real execution may clarify whether malformed JSON should remain `INVALID_SNAPSHOT_SCHEMA` or become a future `INVALID_JSON` code.
- The official conformance pack now includes a representable unsupported canonicalization profile case for `UNSUPPORTED_CANONICALIZATION_PROFILE`.
- Real execution may require additional representable cases for unsupported reconstruction profiles and missing references before stable promotion.

Stable-readiness checklist
--------------------------

Before promoting Core 2.0 reason codes to stable:

1. Confirm the registry, Verify Result v2 draft, JSON schema enum, examples, and conformance fixtures contain the same reason-code set.
2. Confirm every code has a clear `result` group and non-overlapping meaning.
3. Confirm partner feedback has not identified missing high-priority diagnostic codes.
4. Confirm runner real execution does not expose ambiguous or overloaded reason-code meanings.
5. Confirm stable release notes state that reason-code identifiers are compatibility-sensitive.
6. Confirm future additions follow an additive, non-renaming compatibility policy.

Change history
--------------

- 2026-06-02 — Added reason-code review and stabilization notes for Core 2.0 draft stabilization.
- 2026-06-02 — Added runner real execution implications for future reason-code validation.
- 2026-06-02 — Documented current malformed JSON mapping to `INVALID_SNAPSHOT_SCHEMA`.
