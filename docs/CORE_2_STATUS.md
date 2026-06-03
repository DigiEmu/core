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
 - Core 2.0 runner real execution design (draft / not implemented)
 - Core 2.0 runner Phase 2 input JSON parsing (malformed JSON classification)
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
 - Core 2.0 Draft 3 readiness notes (docs/CORE_2_DRAFT_3_READINESS.md)

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

- Runner execution design: `docs/CORE_2_RUNNER_EXECUTION_DESIGN.md` defines a
	post-`v2.0.0-draft.2` path from structural expected-result validation toward
	real `input.json` parsing, verify execution, and actual-vs-expected comparison
	without changing current runner behavior.

- Runner input parsing: the internal conformance runner now reads `input.json`
	and deterministically maps malformed JSON to `ERROR` /
	`INVALID_SNAPSHOT_SCHEMA` for expected-result comparison while preserving
	representative expected-result behavior for well-formed inputs.

- Draft 3 readiness notes: `docs/CORE_2_DRAFT_3_READINESS.md` summarizes
	planned `v2.0.0-draft.3` milestone readiness, required pre-tag checks, and
	non-goals without creating tags or promoting Core 2.0 to stable.

- Actual-vs-expected comparison design: `docs/CORE_2_ACTUAL_EXPECTED_COMPARISON_DESIGN.md`
	defines the next post-draft.3 runner hardening step for comparing observed
	verify results with `expected_verify_result.json` before implementation.

- Runner observed-vs-expected comparison helper: the internal conformance runner
	now compares observed and expected `result` / `reason_code` values through a
	dedicated helper while keeping the public JSON report shape unchanged.

- Runner comparison integration: malformed input JSON observed results now flow
	through the centralized observed-vs-expected comparison path, preserving
	`ERROR` / `INVALID_SNAPSHOT_SCHEMA` behavior.

- Unsupported profile coverage: the official Core 2.0 conformance pack now has
	11 cases, including an unsupported canonicalization profile case represented
	as `ERROR` / `UNSUPPORTED_CANONICALIZATION_PROFILE`.

- Post-draft.3 hardening summary: `docs/CORE_2_POST_DRAFT_3_HARDENING_SUMMARY.md`
	summarizes the completed runner/conformance hardening checkpoint and records
	the current 11-case conformance status (`total=11 passed=11 failed=0`).

- Missing-reference coverage: `MISSING_REFERENCE` is not currently supported by
	the Verify Result v2 schema reason-code enum, so it remains pending and the
	official conformance pack remains at 11 cases. Schema consideration is
	documented in `docs/CORE_2_MISSING_REFERENCE_SCHEMA_CONSIDERATION.md`.

- Draft 4 readiness planning: `docs/CORE_2_DRAFT_4_READINESS.md` prepares a
	future `v2.0.0-draft.4` milestone without tagging it and records the current
	11-case conformance status (`total=11 passed=11 failed=0`).

- Draft 4 post-tag note: `docs/CORE_2_DRAFT_4_POST_TAG_NOTE.md` records that
	main includes a post-`v2.0.0-draft.4` review-fix for the basic conformance
	report fixture without moving or recreating the draft tag.

- Draft 4 self-audit helper: `scripts/audit_core2_draft4.ps1` provides a local
	read-only audit for the current Draft 4 / post-tag state, including tag
	presence, clean tree, tests, conformance output, fixture consistency, and
	documentation guardrails.

## Known Risks

- Verify Result v2 remains draft and is not yet active CLI behavior.
- Full runner verify execution remains unimplemented; current hardening is limited to input JSON parsing and internal observed-vs-expected comparison scaffolding.
- Post-quantum support is migration-readiness only; no PQC primitives are enabled.
- Unicode normalization is documented but not changed; normalization policy remains an open decision.
- `json.RawMessage` handling is documented and tests lock current behavior, but it may require a future profile decision.
- Core 2.0 conformance pack is draft and not yet an executable certification suite.
- Reason-code identifiers are draft but compatibility-sensitive for partner automation.

## Next Recommended Steps

1. Design and implement runner actual-vs-expected comparison while keeping the public JSON report shape unchanged.
2. Extend runner Phase 2 from JSON parsing into deterministic input shape/schema classification.
3. Keep the reason-code registry, Verify Result v2 draft, schema enum, examples, and conformance fixtures aligned.
4. Decide `json.RawMessage` handling for future canonicalization profiles and document migration steps.
5. Collect partner feedback on reason-code clarity before stable promotion.

## Compatibility Statement

All Core 2.0 drafts and tests are designed to be forward-compatible and to avoid changing the v1.0 CLI contract or production behavior. The repository currently contains only draft artifacts and tests that validate and lock existing canonicalization and hashing semantics; no runtime code was modified.

---

If you want, I can open a follow-up PR that adds the conformance runner and extends the conformance vectors — otherwise, the current branch is ready to commit the documentation snapshot.

Note: The repository README now highlights the Core 2.0 Draft 1 partner-testable milestone
and links to the conformance quickstart and release checklist (docs/CORE_2_RELEASE_CHECKLIST.md).
