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
- Reason code review and stabilization notes
- Verify Result v2 JSON Schema draft
- Verify Result v2 example validation test
- Core 2.0 Conformance Pack (draft)
 - Core 2.0 Conformance Pack (draft)
 - Expanded conformance pack (10 cases) including additional negative and edge-case vectors (testdata/core_2_conformance)
- Canonicalization audit
- Canonicalization behavior tests
 - Profile Registry (draft)
 - Hash boundary vectors and tests
 - Canonicalization decision record (draft)
 - Internal conformance runner MVP (draft)
 - Internal conformance runner MVP (draft)
 - Experimental conformance CLI (draft)
 - Conformance quickstart (partner-facing)
 - Core 2.0 versioning guidance (partner-facing)
 - Core 2.0 compatibility matrix (partner-facing)
 - Core 2.0 partner integration notes (partner-facing)
- Core 2.0 release notes (draft)
- Core 2.0 roadmap (next)
- JSON conformance report schema (draft/schema-backed)
 - Core 2.0 release checklist (docs/CORE_2_RELEASE_CHECKLIST.md)
 - Core 2.0 tagging plan (docs/CORE_2_TAGGING_PLAN.md)

## OpenAPI contract draft

- `docs/CORE_2_OPENAPI_DRAFT.md` — OpenAPI contract draft describing intended HTTP surfaces for future partner integrations (documentation-only). This is a draft and does not implement or expose any network API in the repository.
- `openapi/core_2_conformance_api.yaml` — Draft OpenAPI 3.x specification for the above.
 - `docs/CORE_2_DOCKER_USAGE.md` — Optional Docker-based partner usage path for running the experimental conformance CLI without installing Go.

## Draft Artifacts

- `docs/CORE_2_VERIFY_RESULT_v2_DRAFT.md` — v2 draft and rationale
- `schemas/verify_result_v2.schema.json` — v2 JSON Schema (draft)
- `docs/CORE_2_REASON_CODES.md` — proposed reason-code registry
- `docs/CORE_2_REASON_CODE_REVIEW.md` — reason-code review and stabilization notes
- `docs/CORE_2_CONFORMANCE_PACK.md` and `testdata/core_2_conformance/` — draft conformance vectors
- `docs/CORE_2_CANONICALIZATION_AUDIT.md` — audit and recommendations
 - `docs/CORE_2_PROFILE_REGISTRY.md` — Core 2.0 profile registry (draft)
 - `schemas/core_2_profile_registry.schema.json` — profile registry JSON Schema (draft)
 - `testdata/core_2_profiles/profile_registry_valid.json` — example registry

## Tested Guarantees

- Verify Result v2 examples validate against the draft schema.
- Canonical JSON behavior is locked by tests.
- Snapshot Hash v1 mutation behavior is tested.
- Replay outside-hash behavior is tested.
 - OpenAPI draft YAML parses and contains required paths (structural test).

## CI Conformance Checks

- A GitHub Actions CI job now runs `go test ./...` and the experimental Core 2.0 conformance
	CLI (human-readable and `--json`) against the `testdata/core_2_conformance` pack to
	exercise the draft conformance vectors automatically on pull requests.

- A `docker-conformance` CI job builds the repository Docker image and runs the
  experimental conformance CLI inside the container (human-readable + `--json`) to
  validate the Docker usage path.

- Partner handoff package: `docs/CORE_2_PARTNER_HANDOFF.md` added to provide a
	concise partner-facing checklist and review guidance for draft tags like
	`v2.0.0-draft.2`.

- Feedback intake: partner feedback templates and a short feedback process
	(`docs/CORE_2_FEEDBACK_PROCESS.md`) were added to streamline issue intake
	and triage for conformance feedback, integration blockers, and spec
	clarifications.

- Draft 2 review summary: `docs/CORE_2_DRAFT_2_REVIEW_SUMMARY.md` documents the
	automated/Copilot review outcome, confirmed strengths, out-of-scope items,
	and recommended next steps for `v2.0.0-draft.2`.

- Reason-code stabilization review: `docs/CORE_2_REASON_CODE_REVIEW.md` records
	the current draft reason-code set, release risks, and stable-readiness
	checklist without changing schemas, testdata, CLI behavior, or production code.

## Known Risks

- Verify Result v2 remains draft and is not yet active CLI behavior.
- Post-quantum support is migration-readiness only; no PQC primitives are enabled.
- Unicode normalization is documented but not changed; normalization policy remains an open decision.
- `json.RawMessage` handling is documented and tests lock current behavior, but it may require a future profile decision.
- Core 2.0 conformance pack is draft and not yet an executable certification suite.
- Reason-code identifiers are draft but compatibility-sensitive for partner automation.

## Next Recommended Steps

1. Keep the reason-code registry, Verify Result v2 draft, schema enum, examples, and conformance fixtures aligned.
2. Decide `json.RawMessage` handling for future canonicalization profiles and document migration steps.
3. Expand inside/outside-hash conformance vectors to cover more real-world fixtures.
4. Collect partner feedback on reason-code clarity before stable promotion.
5. Plan secure-layer signature work and post-quantum migration as separate engineering projects.

## Compatibility Statement

All Core 2.0 drafts and tests are designed to be forward-compatible and to avoid changing the v1.0 CLI contract or production behavior. The repository currently contains only draft artifacts and tests that validate and lock existing canonicalization and hashing semantics; no runtime code was modified.

---

If you want, I can open a follow-up PR that adds the conformance runner and extends the conformance vectors — otherwise, the current branch is ready to commit the documentation snapshot.

Note: The repository README now highlights the Core 2.0 Draft 1 partner-testable milestone
and links to the conformance quickstart and release checklist (docs/CORE_2_RELEASE_CHECKLIST.md).
