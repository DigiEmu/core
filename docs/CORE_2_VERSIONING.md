# Core 2.0 — Versioning and Compatibility Guidance

Purpose
-------

This document describes the recommended versioning, compatibility, and stability rules for Core 2.0 draft artifacts, profiles, schemas, and experimental tooling produced in this repository.

Principles
----------

- v1.0 behavior remains stable: implementation and CLI compatibility with the v1.0 contract must not be broken by Core 2.0 draft work.
- Core 2.0 artifacts are draft unless explicitly marked stable. Draft artifacts are informative and intended for feedback and testing.
- Experimental CLI commands MUST remain under the `experimental` namespace to avoid accidental promotion to stable surfaces.
- Schema versions must be explicit and included in any schema-backed artifact (for example: `verify_result_v2.schema.json`, profile registry versions, etc.).
- Profile identifiers must be explicit and registered in the profile registry to enable repeatable canonicalization/hashing decisions.
- Changes to canonicalization behavior require a new profile and MUST be documented in the profile registry and canonicalization decision records.
- Any breaking change to the public verification model requires a new major profile or schema version.

Version Labels and Meaning
--------------------------

- `stable-v1` — The existing v1.0 behavior and artifacts. Production implementations that rely on v1.0 must continue to be supported.
- `draft-core-2` — Draft artifacts for Core 2.0. These are experimental and subject to change; they are intended for partner testing and feedback.
- `experimental` — Tooling or CLI surfaces that are intentionally non-stable and should be invoked under `digiemu experimental ...`.
- `schema-backed` — Artifacts that are validated via JSON Schema. Schema versions must be present in filenames and metadata.
- `tested` — Artifacts that have automated tests in this repository to lock behavior.
- `planned` — Features or components that are proposed but not implemented.

Profiles and Canonicalization
-----------------------------

- Canonicalization profiles MUST be explicitly named and included in the Registry. Each profile binds a set of canonicalization, normalization, and hashing decisions.
- When a profile changes canonicalization rules in a way that is not backwards-compatible, a new profile identifier must be created (for example `core-2-profile-v1`, `core-2-profile-v2`).
- Existing snapshots and hashes created under `Snapshot Hash v1` remain valid and must not be reinterpreted under a new canonicalization profile.

Compatibility Guidance
----------------------

- Use the compatibility matrix to determine which artifacts are safe to use together.
- Prefer `stable-v1` artifacts for production workloads.
- Use `draft-core-2` artifacts for testing and feedback; do not rely on drafts for production security guarantees.

Release and Promotion
---------------------

- Promotion from `draft-core-2` to `stable` requires formal review, schema stability, and explicit documentation changes.
- Experimental CLI surfaces must be ported to stable only after a clear migration path and compatibility guarantees are provided.

Appendix: Quick rules
---------------------

1. Never change v1.0 behavior in code that affects runtime contracts.
2. Always include version/type metadata in schema-backed artifacts.
3. Document canonicalization decisions per profile in `docs/CORE_2_CANONICALIZATION_DECISIONS.md`.
4. Treat `testdata/` as a source of reference test vectors; do not change them without a clear test update plan.
