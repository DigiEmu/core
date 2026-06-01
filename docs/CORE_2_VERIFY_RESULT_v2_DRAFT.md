# DigiEmu Core — Verify Result v2 (Draft)

Status: DRAFT

Purpose: Draft of `Verify Result` v2 structure and semantics. This draft is intended to be backward-compatible and optionally emitted by Core 2.0 in an opt-in mode.

High-level goals
----------------
- Provide clearer machine-readable fields for failures, evidence references, and cryptographic profile metadata.
- Preserve v1.0 outputs as-is; v2 SHALL be opt-in and non-destructive.

Result values
-------------
Allowed values for the top-level `result` field are: `PASS`, `FAIL`, `ERROR`.

Required fields (MUST)
- `verify_result_version` (string) — e.g. `v2-draft`
- `result` — one of `PASS` | `FAIL` | `ERROR`
- `snapshot_hash` — canonical snapshot hash computed inside the Core boundary
- `canonicalization_profile` — declared profile used for canonicalization
- `hash_algorithm` — identifier of the hash algorithm used (e.g. `sha256`)
- `reconstruction_profile` — reference to reconstruction rules/profile used
- `verifier_version` — Core or verifier implementation version
- `reason_code` — one of the standardized reason codes (see Reason Codes section)

Optional fields (MAY)
- `snapshot_hash_expected` — expected/declared snapshot hash when a comparison occurs
- `snapshot_hash_actual` — actual computed snapshot hash
- `failure_path` — dotpath or JSON Pointer to the failing field or area
- `evidence_references` — list of references to external evidence (storage URIs, content-addressed IDs)
- `signature_status` — high-level status of any attached signatures (e.g., `none` / `detached` / `embedded`)
- `crypto_profile` — optional reference to cryptographic profile used for signing or multi-hash
- `outside_hash_metadata_status` — validation status for metadata kept outside the snapshot hash

Reason codes (canonical list)
-----------------------------
PASS
- `STATE_RECONSTRUCTED`

FAIL
- `HASH_MISMATCH`
- `RECONSTRUCTION_RULE_MISMATCH`
- `CHAIN_BROKEN`
- `TRANSITION_INVALID`
- `INSIDE_HASH_FIELD_CHANGED`

ERROR
- `CANONICALIZATION_ERROR`
- `INVALID_SNAPSHOT_SCHEMA`
- `MISSING_REQUIRED_FIELD`
- `UNSUPPORTED_HASH_ALGORITHM`
- `UNSUPPORTED_CANONICALIZATION_PROFILE`
- `UNSUPPORTED_RECONSTRUCTION_PROFILE`
- `INTERNAL_ERROR`

Semantics notes (normative)
---------------------------
- `snapshot_hash` MUST be computed only from the inside-hash fields defined by the canonicalization profile.  
- `snapshot_hash_expected` vs `snapshot_hash_actual`: When present, these fields SHOULD be used by consumers to explain a `HASH_MISMATCH` failure.
- `failure_path` SHOULD be provided for `FAIL` cases where field-level diffs are computable.
- `evidence_references` MAY point to artifacts stored in DigiEmu Secure or other storage systems; Core MUST NOT assume custody of those artifacts.

Compatibility
-------------
Core 2.0 SHALL continue to emit existing v1.0 `Verify Result` outputs. v2 outputs are OPTIONAL and MUST NOT be emitted by default for v1.0 CLI commands unless explicitly enabled.

Registry and schema
-------------------
The canonical reason-code registry is published in `docs/CORE_2_REASON_CODES.md` and a machine-readable JSON Schema draft is available at `schemas/verify_result_v2.schema.json`. Example v2 `Verify Result` documents are provided in `testdata/verify_result_v2/` for PASS/FAIL/ERROR cases.
