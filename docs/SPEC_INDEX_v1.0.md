# DigiEmu Core — SPEC_INDEX_v1.0

**Status:** DRAFT (Phase 5)  
**Purpose:** One canonical index of normative contracts, public review drafts, and implementation-facing specification documents.

---

## Normative contracts

These documents define the current implementation-facing contracts for DigiEmu Core v1.0.

- `docs/CLI_CONTRACT_v1.0.md`
- `docs/VERSIONING_POLICY_v1.0.md`
- `docs/CLI_VERIFY_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/VERIFY_RESULT_SCHEMA_v1.json`

---

## Public review drafts

These documents define the emerging public standard structure for DigiEmu Core.

- `docs/DIGIEMU_CORE_SPEC_v0.9.md`
- `docs/TEST_VECTORS_v0.9.md`
- `docs/NEGATIVE_TEST_VECTORS_v0.9.md`
- `docs/TEST_VECTOR_MANIFEST_v0.9.json`
- `docs/CONFORMANCE_v0.9.md`
- `docs/CONFORMANCE_DECLARATION_v0.9.md`
- `docs/CONFORMANCE_DECLARATION_SCHEMA_v0.9.json`
- `docs/VERIFY_REPORT_EXAMPLES_v0.9.md`
- `docs/VERIFY_REPORT_SCHEMA_v0.9.json`

---

## Supporting specification documents

These documents provide additional specification context for snapshot hashing, bundle layout, verification semantics, and security boundaries.

- `docs/SNAPSHOT_HASH_v1.0.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
- `docs/SECURITY.md`
- `docs/THREAT_MODEL.md`

---

## Public standard structure

DigiEmu Core currently follows this public standard development path:

```text
Specification → Test Vectors → Conformance → Conformance Declaration → Conformance Declaration Schema → Verify Report Examples → Verify Report Schema
```

Meaning:

```text
Specification explains the model.
Test vectors make verification reproducible.
Conformance defines implementer requirements.
Verify report examples define machine-readable outcomes.
Verify report schema makes verification reports formally validatable.
```

---

## Core 2.0 hardening drafts

These documents are forward-compatible Core 2.0 hardening drafts and do not replace or break the current v1.0 contract.

- `docs/CORE_2_BOUNDARY_MODEL.md`
- `docs/CORE_2_HARDENING_PLAN.md`
- `docs/CORE_2_CRYPTO_AGILITY.md`
- `docs/CORE_2_VERIFY_RESULT_v2_DRAFT.md`
- `docs/CORE_2_MIGRATION_FROM_v1.md`
 - `docs/CORE_2_REASON_CODES.md`
 - `schemas/verify_result_v2.schema.json`
 - `docs/CORE_2_CONFORMANCE_PACK.md`
 - `docs/CORE_2_TEST_VECTOR_RULES.md`
 - `testdata/core_2_conformance/`  
 - `docs/CORE_2_STATUS.md`
 - `docs/CORE_2_PROFILE_REGISTRY.md`
 - `schemas/core_2_profile_registry.schema.json`
 - `testdata/core_2_profiles/profile_registry_valid.json`
 - `docs/CORE_2_HASH_BOUNDARY_VECTORS.md`
 - `testdata/core_2_hash_boundary/`
 - `docs/CORE_2_CANONICALIZATION_DECISIONS.md`
 - `docs/CORE_2_CONFORMANCE_RUNNER.md`
 - `docs/CORE_2_CONFORMANCE_QUICKSTART.md`
 - `docs/CORE_2_PARTNER_INTEGRATION_NOTES.md`
 - `docs/CORE_2_VERSIONING.md`
 - `docs/CORE_2_COMPATIBILITY_MATRIX.md`
 - `docs/CORE_2_RELEASE_NOTES_DRAFT.md`
 - `docs/CORE_2_ROADMAP_NEXT.md`
 - `docs/CORE_2_CONFORMANCE_REPORT.md`
 - `schemas/core_2_conformance_report.schema.json`