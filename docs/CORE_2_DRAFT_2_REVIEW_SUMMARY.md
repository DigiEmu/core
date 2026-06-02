# Core 2.0 Draft 2 Review Summary

Reviewed milestone
-----------------

- Tag: `v2.0.0-draft.2`
- Status: draft / pre-release / not stable
- Meaning: Docker-CI validated partner-testable Core 2.0 draft milestone

Review scope
------------

- Core 2.0 conformance path
- Experimental conformance CLI
- JSON conformance report
- Conformance report schema
- OpenAPI contract draft
- Docker usage path
- Docker CI validation
- Partner handoff readiness
- Feedback intake process

Review verdict
--------------

- `v2.0.0-draft.2` is ready for partner handoff as a draft milestone.
- No critical blockers were identified in the automated/Copilot review.
- Core 2.0 remains draft/pre-release and is not promoted to stable by this review.
- `v1.0` behavior remains unchanged and is the compatibility baseline.

Confirmed strengths
-------------------

- Partner-testable conformance path (local and Docker)
- Human-readable conformance output (summary)
- Machine-readable JSON conformance report
- Schema-backed report format and draft schema
- OpenAPI contract draft and structural validation
- Docker-based usage path and CI validation
- CI conformance checks for PRs
- Partner handoff documentation and structured feedback intake
- Reason-code review and stabilization notes added for draft-to-stable readiness
- Runner real execution design identified as the next conformance stabilization step

Known out-of-scope items
------------------------

- Stable Core 2.0 CLI and stable promotion
- Production HTTP API server
- Authentication, tenancy, and policy judgement
- Secure Layer signatures and PQC implementations
- SDKs for non-Go languages
- Legal certification or formal accreditation

Recommended next PRs
--------------------

High priority
- Conformance pack expansion (more cases) — completed in this branch (added 7 edge/negative cases; total 10)
- Additional negative and edge-case test vectors
- Reason code clarity and coverage review — documented in `docs/CORE_2_REASON_CODE_REVIEW.md`
- Runner real execution design — documented in `docs/CORE_2_RUNNER_EXECUTION_DESIGN.md`

Medium priority
- OpenAPI contract refinement and examples
- Partner feedback triage and first-response plan
- Docker usage refinement (size/runtime ergonomics)

Later
- API Server MVP (if partners request remote validation)
- Secure Layer Signature MVP
- Post-Quantum migration profile planning
- SDK and partner example planning

Important notes
---------------

- This review summary does not promote Core 2.0 to stable.
- Experimental commands remain under the `experimental` namespace.
- Draft artifacts may still evolve before a stable `v2.0.0` release.
- Reason-code identifiers are draft but should be treated as compatibility-sensitive before stable promotion.
- Runner real execution remains design-only and is not implemented by this review.
- Partner feedback should be submitted using the issue templates in
  `.github/ISSUE_TEMPLATE/` and the process described in
  `docs/CORE_2_FEEDBACK_PROCESS.md`.

Change history
--------------

- 2026-06-02 — Initial review summary added for `v2.0.0-draft.2`.
- 2026-06-02 — Added reason-code review and stabilization-note reference.
- 2026-06-02 — Added runner real execution design reference as the next stabilization step.
