# DigiEmu Core — VERIFY_SPEC_v1.0

Version: v1.0
Status: DRAFT
Scope: Verify Operator + ReplayV1 + canonical_json_v1 hashing + write-policy

## 1. Purpose

This document defines, normatively, how a bundle SHALL be deterministically reconstructed into a replayed state and verified.

Auditability requirement:

- Each verification output SHALL be reconstructable from its inputs (bundle files and CLI inputs) under the deterministic replay and hashing rules defined here.

## 2. Terms

The following terms are used in this document:

- **bundle**: A set of files forming the verification input for a single snapshot reference.
- **bundle root**: The directory `snapshots/<ref>/` that contains `snapshot.json` and optional subdirectories.
- **snapshot.json**: The required JSON file located at `<bundleRoot>/snapshot.json`.
- **claim file**: A JSON file located at `<bundleRoot>/claims/*.json`.
- **replayed state**: The deterministic, reconstructed JSON object produced by ReplayV1.
- **canonical_json_v1**: A stable canonical JSON encoding used for hashing.
- **canonical_scope**: A string identifier describing what is included in the canonicalized bytes.
- **expected_hash_v1**: The declared expected digest stored in `snapshot.json`.
- **trace**: An ordered list of strings that records which inputs were attempted and read.

## 3. Inputs

### 3.1 Bundle Root Layout

The bundle root SHALL have the following logical layout:

- `snapshots/<ref>/`
	- `snapshot.json` (required)
	- `claims/*.json` (optional; if present: MUST be included in ReplayV1)

Implementations MAY support additional optional directories (e.g. `units/`, `versions/`, `meaning/`, `uncertainty/`), provided determinism requirements are met.

### 3.2 CLI inputs

The CLI SHALL accept the following inputs:

- `--ref` (snapshot reference)
- `--bundle` (path to `snapshots/<ref>`)
- `--data`, `--fixture-root`, `--prefer-data` (bundle lookup rules)

Lookup rules:

- If `--bundle` is provided, the resolved bundle root SHALL be that path.
- Otherwise, the CLI SHALL resolve the bundle root by searching `<fixture-root>/snapshots/<ref>/` and `<data>/snapshots/<ref>/`.
- If `--prefer-data` is set, the CLI SHALL attempt the data path before the fixture path.

## 4. Deterministic Replay (ReplayV1)

ReplayV1 MUST reconstruct the state deterministically from bundle inputs.

ReplayV1 MUST NOT compute the hash.

ReplayV1 MUST return (or otherwise provide) a trace list that captures key input files and the used bundle root.

Trace MUST include:

- the resolved bundle root path entry (e.g. `"used:<bundleRoot>"`)
- the `snapshot.json` path
- each claim file path that was loaded

Replay inclusion order:

- For `claims/*.json`, ReplayV1 MUST include files in lexicographic filename order.
- Only files with a `.json` suffix (case-insensitive) SHALL be included.

## 5. Canonicalization + Hashing

hash_alg MUST be `"sha256(canonical_json_v1)"`.

canonical_json_v1 MUST be a stable, canonical JSON encoding.

canonical_scope MUST be `"canonical_json_v1_excluding_expected_hash_v1"`.

expected_hash_v1 MUST be excluded from the hashed scope to prevent self-referential hashing.

Normatively:

- Before hashing, the operator MUST remove the `expected_hash_v1` key from the `snapshot` content included in the replayed state.

Pretty printing or output formatting MUST NOT affect the hash.

## 6. Verify Procedure (Operator.Verify)

The operator SHALL implement the following sequence:

1) Resolve bundle root (via `--bundle` or lookup using `--data`/`--fixture-root`/`--prefer-data`).
2) Load bundle inputs (BOM-safe JSON decode).
3) ReplayV1(bundle) -> replayed state + trace.
4) Compute `got = HashV1FromState(replayed state)` (or equivalent).
5) Compare `got` vs `expected` (read from `snapshot.json.expected_hash_v1`).
6) Produce ResultV1 (JSON contract) including `trace`.
7) Optionally apply write-policy if `--write-expected` is set.

## 7. Result Contract (JSON)

The operator SHALL emit a single JSON object containing the following fields:

- `ok` (bool): true if verification succeeded.
- `ref` (string): the snapshot reference.
- `expected` (string): the declared `expected_hash_v1`.
- `got` (string): the computed digest.
- `hash_alg` (string): hashing algorithm identifier.
- `canonical_scope` (string): canonicalization scope identifier.
- `trace` (string[]): ordered trace entries.
- `message` (string): human-readable status message.
- `wrote_expected` (bool): true if the operator wrote `expected_hash_v1`.
- `write_blocked` (bool): true if a write was attempted/considered but blocked.
- `write_reason` (enum string): reason code for write outcome.

write_reason allowed values MUST be exactly one of:

- `"none"`
- `"flag_not_set"`
- `"placeholder"`
- `"existing_expected_present"`
- `"snapshot_not_found"`
- `"snapshot_invalid_json"`
- `"io_error"`
- `"invalid_hash"`

## 8. Write Policy (--write-expected)

The operator MUST NOT overwrite a non-placeholder `expected_hash_v1`.

The operator MAY write `expected_hash_v1` only if the current value is a recognized placeholder as defined by the placeholder policy function name `IsPlaceholderExpected`.

On blocked write due to existing expected:

- `write_blocked` MUST be true
- `write_reason` MUST be `"existing_expected_present"`

Writer requirements:

- Writer MUST write JSON without a UTF-8 BOM.
- Writer MUST preserve stable/pretty JSON formatting (indentation is allowed; output MUST remain valid JSON).

## 9. Exit Codes

The CLI SHALL map outcomes to exit codes as follows:

- `0`: ok (verify ok; and no governance-block event)
- `2`: mismatch (`got != expected`) when treated as a verification failure
- `3`: write blocked due to existing expected present, when `--write-expected` is set
- `4`: invalid input (missing snapshot, invalid json, invalid hash, invalid bundle path, etc.)
- `5`: internal error

Verification success (`ok=true`) MAY still return exit 3 if `--write-expected` is set and the write was blocked due to existing expected.

## 10. Determinism Requirements

- BOM-safe decode MUST be applied.
- Canonical JSON MUST be stable for the same logical content.
- Trace ordering MUST be deterministic.

## 11. Security & Governance Notes

No-overwrite policy rationale:

- The operator MUST be safe against accidental overwrites of `expected_hash_v1` caused by human error.

Recommended CI policy:

- treat exit 3 as a governance violation requiring review
- treat exit 2 as verification failure
- treat exit 4 and 5 as hard failures

Operational recommendation:

- keep fixtures tracked (version controlled)
- keep runtime data untracked (environment-specific)

## 12. References

- `docs/CLI_VERIFY_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
- `docs/SNAPSHOT_HASH_v1.0.md`
- internal packages/types (by name only): `internal/verify.Operator`, `internal/verify.ReplayV1`, `pkg/snapshot.HashV1FromState`
