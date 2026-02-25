# Phase 3.0 — Task 1/8: Domain Types (Claim / ClaimSet / Relations)

Spec ID: `p3-0-t1-domain-types`
Phase: 3.0
Status: done
Kind: spec/runbook

Purpose
-------
Document the implemented domain types for Claim, ClaimSet and ClaimRelation and the minimal validation (`ValidateMinimal`) added in this task.

Constraints and Goals
---------------------
- Add minimal first-class domain types for Claims.
- Keep everything additive and optional.
- No breaking changes; deterministic-ready; no storage logic in domain.

Repository locations
--------------------
- Domain directory: internal/kernel/domain
- Ports directory: internal/kernel/ports

Implemented additions
---------------------
- `internal/kernel/domain/claim.go`
  - Types: `Claim`, `ClaimSet`, `ClaimRelation`, `ClaimRelationType`
  - Constants: `ClaimSchemaV0`, `ClaimSetSchemaV0`
  - Method: `ValidateMinimal()` on `ClaimSet` enforcing the minimal rules below.

- `internal/kernel/domain/claim_test.go`
  - Unit tests covering positive and negative validation cases.

Schema versions
---------------
- Claim schema: `claim/v0`
- ClaimSet schema: `claimset/v0`

Shape (summary)
----------------
- Claim
  - `schema_version` (string, required, must equal `claim/v0`)
  - `id` (string, required)
  - `text` (string, required)
  - `tags` ([]string, optional)

- ClaimSet
  - `schema_version` (string, required, must equal `claimset/v0`)
  - `version_id` (string, required)
  - `claims` ([]Claim, required)
  - `relations` ([]ClaimRelation, optional)

- ClaimRelation
  - `type` (enum: `CONTRADICTS`)
  - `from_claim_id` (string)
  - `to_claim_id` (string)

Validation rules (ValidateMinimal)
----------------------------------
- Claim.schema_version must equal `claim/v0`.
- ClaimSet.schema_version must equal `claimset/v0`.
- Claim.id and Claim.text must be non-empty.
- ClaimSet.version_id must be non-empty.
- All Claim IDs in a ClaimSet must be unique.
- Each relation must reference existing claim ids.
- Relation type must be supported (`CONTRADICTS`).

Permissive choices (out of scope)
---------------------------------
- No semantic validation of claim text.
- No ontology/tag validation.
- No uncertainty model included here.

Acceptance criteria
-------------------
- New domain types compile and do not affect existing packages.
- No changes required to existing meaning-layer code.
- `ValidateMinimal()` catches empty ids/text and invalid relation references.
- All existing tests remain green.

Tests added
-----------
- `internal/kernel/domain/claim_test.go`
  - Cases:
    - valid ClaimSet passes `ValidateMinimal()`
    - duplicate claim ids fails
    - relation referencing missing id fails
    - unknown relation type fails
    - wrong schema_version fails

Verified
--------
- Ran `gofmt -w .` and `go test ./...` — all tests passed in the repository after implementing these files.

Notes / Next tasks
------------------
- This task is intentionally additive and prepares for hashing, sidecars and audit events in following tasks.
- Out-of-scope for this task: hashing helpers, audit events, adapters, verify-audit, export, CLI/HTTP wiring, docs/release.
