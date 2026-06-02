# Core 2.0 Partner Handoff — v2.0.0-draft.2

Status: draft / pre-release / not stable

This document provides a concise partner-facing handoff for the `v2.0.0-draft.2`
milestone. It explains what partners should test, review, and provide feedback on
when evaluating the current Core 2.0 draft artifacts.

Current milestone
-----------------

- Tag example: `v2.0.0-draft.2`
- Status: draft / pre-release / not stable
- Meaning: Docker-CI validated partner-testable Core 2.0 draft milestone

Who this is for
----------------

- Technical partners
- AI governance reviewers
- Conformance implementers
- Integration partners
- External reviewers

What to test
------------

- Run the Go test suite: `go test ./...`
- Run the experimental conformance CLI (human-readable):
  `go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance`
- Run the experimental conformance CLI (JSON):
  `go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance --json`
- Validate the Docker-based conformance path (optional):
  `docker build -t digiemu-core .` and run the two `docker run` commands.
- Review the OpenAPI contract draft (`openapi/core_2_conformance_api.yaml`).

Commands (local)

```bash
git checkout v2.0.0-draft.2
go test ./...
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance --json
```

Commands (docker)

```bash
docker build -t digiemu-core .
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance --json
```

Expected results
----------------

Human:
- `Conformance run summary: total=3 passed=3 failed=0`

JSON:
- `report_version: core-2-conformance-report-v1`
- `status: PASS`
- `total: 3`
- `passed: 3`
- `failed: 0`

What to review
--------------

- Are the conformance cases understandable and actionable?
- Are reason codes clear and sufficient for diagnosing failures?
- Is the JSON report shape suitable for automation and ingestion?
- Does the OpenAPI contract draft capture the needed endpoints and data shapes?
- Is the Docker usage path practical for partners without Go tooling?
- Are the boundaries between draft and stable clear enough for partner adoption?

Feedback requested
------------------

- Missing conformance edge cases or gaps in coverage
- Ambiguous or missing reason codes
- Additional report fields required by partner automation
- OpenAPI contract shape or endpoint suggestions
- Docker usage friction or environment constraints
- Integration blockers or compatibility concerns

Out of scope
------------

- Stable Core 2.0 release (this is a draft)
- Production HTTP API server implementation
- Authentication, tenancy, or policy judgement
- Secure Layer signatures or PQC implementations
- Legal certification or external accreditation

Important notes
---------------

- `v1.0` behavior remains unchanged and is the baseline for compatibility.
- Core 2.0 remains draft unless explicitly marked stable.
- Experimental commands remain under the `experimental` namespace.
- This handoff is for partner review and testing, not production certification.

Links
-----

- Conformance Quickstart: `docs/CORE_2_CONFORMANCE_QUICKSTART.md`
- Partner Integration Notes: `docs/CORE_2_PARTNER_INTEGRATION_NOTES.md`
- Docker Usage: `docs/CORE_2_DOCKER_USAGE.md`
- OpenAPI Draft: `docs/CORE_2_OPENAPI_DRAFT.md`
- Release Checklist: `docs/CORE_2_RELEASE_CHECKLIST.md`
- Tagging Plan: `docs/CORE_2_TAGGING_PLAN.md`
- OpenAPI spec: `openapi/core_2_conformance_api.yaml`
- Conformance report schema: `schemas/core_2_conformance_report.schema.json`

Feedback and review
-------------------

See also: `docs/CORE_2_DRAFT_2_REVIEW_SUMMARY.md` for a short review outcome and
recommended next steps for `v2.0.0-draft.2`.

If you have feedback, please use the partner feedback issue templates in
`.github/ISSUE_TEMPLATE/` and consult `docs/CORE_2_FEEDBACK_PROCESS.md` for
guidance on what to include and how issues will be triaged.

Change history
--------------

- 2026-06-02 — Initial partner handoff draft for v2.0.0-draft.2
