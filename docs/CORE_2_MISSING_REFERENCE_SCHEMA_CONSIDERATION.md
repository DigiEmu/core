# Core 2.0 MISSING_REFERENCE Schema Consideration

**Milestone:** post-v2.0.0-draft.3  
**Stability:** discussion / draft / not implemented  
**Current decision:** `MISSING_REFERENCE` is not currently part of the Verify Result v2 `reason_code` enum.

This document records a documentation-only consideration for whether and how `MISSING_REFERENCE` could be added to Verify Result v2 later. It does not change production code, tests, schemas, testdata, CI behavior, CLI behavior, report shape, tags, or v1.0 compatibility.

## Current State

- `MISSING_REFERENCE` is documented as pending.
- No official conformance case exists for `MISSING_REFERENCE`.
- The official conformance pack remains at 11 cases.
- The public JSON report shape remains unchanged.
- Verify Result v2 schema is unchanged.

## Possible Semantics

If added later, `MISSING_REFERENCE` should have narrow and deterministic semantics:

- A required referenced snapshot, bundle, transition, or external verification input is absent.
- The system cannot complete reconstruction because a declared reference cannot be resolved.
- The failure is not a hash mismatch, not malformed JSON, and not an unsupported profile.
- The error is deterministic and should be reproducible from the input package.

## Boundaries

- `MISSING_REFERENCE` should not be used for malformed JSON.
- `MISSING_REFERENCE` should not be used for missing `input.json` or missing `expected_verify_result.json` in conformance case structure.
- `MISSING_REFERENCE` should not replace `CHAIN_BROKEN` if the chain exists but fails continuity rules.
- `MISSING_REFERENCE` should not replace `RECONSTRUCTION_RULE_MISMATCH`.

## Requirements Before Schema Addition

- Define exact reference types covered.
- Define whether missing references belong to Verify Result v2 or a higher-level bundle validation result.
- Add schema enum value only after semantics are stable.
- Add official conformance case only after schema support exists.
- Add runner tests for observed `ERROR` / `MISSING_REFERENCE`.
- Keep `report_version` unchanged unless report shape changes.

## Risks

- Too broad a reason code could overlap with `CHAIN_BROKEN` or `RECONSTRUCTION_RULE_MISMATCH`.
- Adding the enum too early could create unstable partner contracts.
- Reference resolution may belong outside the current minimal runner scope.

## Recommended Next Steps

- Discuss whether missing references are part of Verify Result v2 or bundle-level validation.
- If accepted, add `MISSING_REFERENCE` to schema in a dedicated PR.
- Then add a representable conformance case.
- Then update CI totals and documentation.

## Compatibility Statement

This consideration does not implement `MISSING_REFERENCE`. Core 2.0 remains draft and pre-release, the current conformance pack remains at 11 cases, and the latest stable line remains `v1.0.0`.

1.0.0 remains the stable/latest release; Core 2.0 remains draft/pre-release.

1.0.0 remains stable/latest; Core 2.0 remains draft/pre-release.

v1.0.0 remains stable/latest; Core 2.0 remains draft/pre-release.
