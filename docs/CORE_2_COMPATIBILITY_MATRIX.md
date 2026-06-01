# Core 2.0 — Compatibility Matrix

This compatibility matrix documents the current status of Core 2.0 components: whether they are stable, draft, experimental, schema-backed, tested, or planned.

Legend
------

- `stable-v1` — Production-stable v1 behavior
- `draft-core-2` — Draft artifact for Core 2.0
- `experimental` — Experimental tooling/CLI under `digiemu experimental`
- `schema-backed` — Validated with a JSON Schema file
- `tested` — Has automated tests in this repo
- `planned` — Intended future feature

Component | Status | Notes
---|---:|---
Snapshot Hash v1 | stable-v1, tested | Existing canonical snapshot hashing; covered by hash boundary vectors
Canonical JSON v1 | stable-v1, tested | Deterministic canonicalization used for HashV1
Verify Result v2 draft | draft-core-2, schema-backed | Draft verify result schema under `schemas/verify_result_v2.schema.json`
Reason Codes Registry | draft-core-2, tested | Draft registry documented in `docs/CORE_2_REASON_CODES.md`
Profile Registry | draft-core-2, schema-backed | Draft schema in `schemas/core_2_profile_registry.schema.json`
Conformance Pack | draft-core-2, tested | Pack under `testdata/core_2_conformance/` for partner testing
Conformance Runner MVP | draft-core-2, tested | Internal runner in `internal/conformance`
Experimental Conformance CLI | experimental | Accessed via `digiemu experimental conformance run ...`
Hash Boundary Vectors | draft-core-2, tested | Vectors under `testdata/core_2_hash_boundary/`
Canonicalization Decision Record | draft-core-2 | Documented in `docs/CORE_2_CANONICALIZATION_DECISIONS.md`
Post-Quantum readiness placeholders | planned | PQC migration is planned; no PQC primitives implemented
Secure Layer signatures | planned | Signature work planned as separate project

Notes
-----

- This matrix is intentionally conservative: draft artifacts should not be treated as production guarantees.
- The `experimental` label indicates tooling that must be invoked explicitly under the `experimental` namespace.
- Schema-backed artifacts include an explicit version identifier in filenames and metadata.
