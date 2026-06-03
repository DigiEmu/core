# Core 2.0 — Roadmap (Next Steps)

This roadmap lists the near-term priorities and planned work items to evolve the
Core 2.0 draft artifacts toward a more comprehensive partner integration and
eventual stabilization path.

Next phase
----------

- Stabilize conformance runner output (CI job added for automated checks)
- Add optional JSON output for the experimental conformance CLI (implemented)
- Add machine-readable conformance report format
- Prepare OpenAPI contract draft for a future HTTP API
 - Prepare OpenAPI contract draft for a future HTTP API (draft created)
- Prepare Docker-based usage path for easier partner evaluation
 - Prepare Docker-based usage path for easier partner evaluation
 - Prepare Docker-based usage path for easier partner evaluation (draft prepared)
 - Docker usage path validated in CI (docker-conformance job)
- Prepare Draft 3 readiness notes as the next milestone preparation step
- Plan Secure Layer signature MVP separately
- Plan Post-Quantum migration profile separately

Priorities
----------

High
- Conformance CLI JSON output
- Conformance report schema
- OpenAPI draft

Medium
- Docker image for partner usage
- Non-Go partner examples and SDK snippets
- Release tagging and clear drafting process
- Draft 3 readiness validation before any future `v2.0.0-draft.3` tag

Later
- Secure Layer Ed25519 signatures MVP
- Hybrid / PQ migration profiles
- Promote a stable Core 2.0 CLI command (after thorough review)

Notes
-----

- The roadmap intentionally keeps experimental tooling under the
  `experimental` namespace until a formal promotion process is defined.
- Work items that affect canonicalization, hashing semantics, or released
  schemas must follow the versioning guidance in `docs/CORE_2_VERSIONING.md`.
