# DigiEmu Core 2.0 — Hardening Status

**Last updated:** 2026-06-01

This document summarizes the current hardening status for DigiEmu Core 2.0: what is completed, what remains a draft, what is covered by tests today, known risks, and recommended next steps. This is a documentation-only status snapshot and does not change any runtime behavior, CLI contract, or schema definitions.

## Current Status

- Development work for Core 2.0 has produced a set of design drafts, validation artifacts, conformance examples, and focused tests that lock current canonicalization and hashing behavior.
- All work is currently staged on the `core-2-status-document` feature branch and is draft-facing: no CLI or production behavior has been changed.

## Completed Hardening Work

- Core 2.0 boundary model
- Core 2.0 hardening plan
- Crypto agility draft
- Migration from v1 guidance
- Verify Result v2 draft
- Reason Codes registry
- Verify Result v2 JSON Schema draft
- Verify Result v2 example validation test
- Core 2.0 Conformance Pack (draft)
- Canonicalization audit
- Canonicalization behavior tests

## Draft Artifacts

- `docs/CORE_2_VERIFY_RESULT_v2_DRAFT.md` — v2 draft and rationale
- `schemas/verify_result_v2.schema.json` — v2 JSON Schema (draft)
- `docs/CORE_2_REASON_CODES.md` — proposed reason-code registry
- `docs/CORE_2_CONFORMANCE_PACK.md` and `testdata/core_2_conformance/` — draft conformance vectors
- `docs/CORE_2_CANONICALIZATION_AUDIT.md` — audit and recommendations

## Tested Guarantees

- Verify Result v2 examples validate against the draft schema.
- Canonical JSON behavior is locked by tests.
- Snapshot Hash v1 mutation behavior is tested.
- Replay outside-hash behavior is tested.

## Known Risks

- Verify Result v2 remains draft and is not yet active CLI behavior.
- Post-quantum support is migration-readiness only; no PQC primitives are enabled.
- Unicode normalization is documented but not changed; normalization policy remains an open decision.
- `json.RawMessage` handling is documented and tests lock current behavior, but it may require a future profile decision.
- Core 2.0 conformance pack is draft and not yet an executable certification suite.

## Next Recommended Steps

1. Define a Core 2.0 profile registry to pin canonicalization and hashing profiles.
2. Add an executable conformance runner or CLI mode to exercise the conformance pack.
3. Decide `json.RawMessage` handling for future canonicalization profiles and document migration steps.
4. Expand inside/outside-hash conformance vectors to cover more real-world fixtures.
5. Plan secure-layer signature work and post-quantum migration as separate engineering projects.

## Compatibility Statement

All Core 2.0 drafts and tests are designed to be forward-compatible and to avoid changing the v1.0 CLI contract or production behavior. The repository currently contains only draft artifacts and tests that validate and lock existing canonicalization and hashing semantics; no runtime code was modified.

---

If you want, I can open a follow-up PR that adds the conformance runner and extends the conformance vectors — otherwise, the current branch is ready to commit the documentation snapshot.
