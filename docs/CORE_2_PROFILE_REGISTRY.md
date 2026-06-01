# DigiEmu Core 2.0 — Profile Registry (Draft)

**Last updated:** 2026-06-01

Purpose
-------

This document provides the draft Core 2.0 Profile Registry: a machine-readable and human-readable registry of canonicalization, hashing, reconstruction, verify-result, conformance, and crypto profiles used by Core 2.0 artifacts and conformance packs. This registry is draft-only and does not change v1.0 runtime behavior or existing snapshot hashes.

Design Principles
-----------------

- The registry is draft-only and does not change v1.0 behavior.
- The registry MUST distinguish `active` profiles from `planned` profiles.
- The registry MUST be machine-readable and versioned.
- The registry MUST support future extension without invalidating v1.0 evidence.
- Post-quantum profiles are migration-readiness placeholders only.

Registry structure (human summary)
----------------------------------

The machine-readable registry (JSON) contains a top-level object with the following named profile categories. Each category contains an `active` array and an optional `planned` array.

- `canonicalization_profiles` — encodings used to canonicalize snapshots for hashing and comparison.
- `hash_algorithms` — digest algorithms used to compute snapshot canonical hashes.
- `reconstruction_profiles` — deterministic reconstruction/replay rules for hashing and verification.
- `verify_result_versions` — supported verify-result schema versions.
- `conformance_levels` — named conformance levels that conformance packs reference.
- `crypto_profiles` — lightweight descriptors for signature/custody models (no key material, no orchestration).

Current registry (human view)
-----------------------------

Active canonicalization profiles:

- digi​emu-canonical-json-v1

Planned canonicalization profiles:

- digi​emu-canonical-json-v2-planned

Active hash algorithms:

- sha256

Planned hash algorithms:

- sha3-256-planned

Active reconstruction profiles:

- digi​emu-reconstruct-v1

Planned reconstruction profiles:

- digi​emu-reconstruct-v2-planned

Active verify-result versions:

- 2.0-draft

Active conformance levels:

- core-2-draft-basic
- core-2-draft-conformance

Active crypto profiles:

- none
- external-signature

Planned crypto profiles (migration placeholders):

- classical-signature-planned
- hybrid-pq-planned
- pq-migration-ready-planned

Notes
-----

- This registry intentionally avoids policy or enforcement mechanisms — it is a registry and reference. Implementation and enforcement (runners, CI, or CLI modes) are separate engineering tasks.
- The machine-readable JSON Schema for the registry is kept in `schemas/core_2_profile_registry.schema.json`.
