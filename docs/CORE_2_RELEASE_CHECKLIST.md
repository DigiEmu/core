# Core 2.0 Draft Release — Checklist

This document defines the minimum set of checks and artifacts that must be present
before creating a reproducible Core 2.0 draft release tag (for example `v2.0.0-draft.1`).

This checklist is documentation-only and does not change code, CLI behavior, or
schemas. It is intended to consolidate the partner-testable milestone into a
clear set of gating criteria for a draft release marker.

Required before creating a draft tag
-----------------------------------

The following MUST be true before creating a draft tag:

- The `main` branch is clean and up-to-date with the intended release content.
- `go test ./...` passes across the repository.
- The experimental conformance CLI human-readable run passes against
  `testdata/core_2_conformance` (shows a summary like `total=3 passed=3 failed=0`).
- The experimental conformance CLI `--json` run passes and emits a machine
  report containing `"status": "PASS"`, `"total": 3`, `"passed": 3`,
  and `"failed": 0` for the included conformance pack.
- CI conformance checks (GitHub Actions) pass for the branch and PR.
- Verify Result v2 schema exists (`schemas/verify_result_v2.schema.json`).
- Conformance report schema exists (`schemas/core_2_conformance_report.schema.json`).
- OpenAPI contract draft exists (`openapi/core_2_conformance_api.yaml`).
- OpenAPI structural validation tests pass (where applicable).
- Docker usage path is documented (`docs/CORE_2_DOCKER_USAGE.md`).
- Partner integration notes are present (`docs/CORE_2_PARTNER_INTEGRATION_NOTES.md`).
- Compatibility matrix exists (`docs/CORE_2_COMPATIBILITY_MATRIX.md`).
- v1.0 behavior remains unchanged and is validated by tests.

Not required for a draft tag
----------------------------

The following are explicitly NOT required to create a draft tag:

- A stable Core 2.0 CLI.
- An HTTP API server implementation.
- Docker image verification in CI.
- Secure Layer signatures or key custody implementations.
- Post-Quantum cryptography implementations.
- SDKs for other languages.
- External certification or formal accreditation.

Usage notes
-----------

- A draft tag (for example `v2.0.0-draft.1`) is intended to be partner-testable,
  reproducible, and to mark a milestone for review and feedback. It is NOT a
  stable or final release.
- Keep experimental commands under the `experimental` namespace; do not
  promote them to stable by tagging alone.

Change history
--------------

- 2026-06-02 — Initial draft checklist added.
