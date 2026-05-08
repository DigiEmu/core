# Snapshot Hash v1.0 — Canonical JSON

Status: NORMATIVE

This document defines the canonical snapshot hashing algorithm used by DigiEmu Core v1.x.

Scope: canonical_utf8_without_sha256_comment_line, SHA-256, uppercase hex.

## Rules

- Build an explicit replay state object that contains only the stable kernel state.
- Serialize using `internal/canonicaljson.Marshal`.
- Object/map keys MUST be sorted deterministically.
- Array order MUST be preserved.
- Insignificant whitespace MUST NOT affect the canonical bytes.
- Compute SHA-256 over the UTF-8 bytes of the canonical JSON.
- Encode the digest as uppercase hexadecimal.

## Deterministic hash boundary

The snapshot hash boundary MUST include only deterministic replay state.

The following values are considered outside the deterministic hash boundary:

- timestamps such as `CreatedAtUnix`
- generated or user-facing labels such as version labels
- actor identifiers such as `ActorID`
- audit trace metadata
- verification metadata such as `expected_hash_v1`
- runtime environment data
- build or bundle metadata that does not change reconstructed state

These values MAY be preserved for auditability, display, traceability, or governance purposes, but they MUST NOT change deterministic state identity unless explicitly included as deterministic replay state by a future versioned specification.

## Replay and verification metadata

`expected_hash_v1` MUST be excluded from the hashed replay scope.

This prevents self-referential hashing and keeps declared verification metadata separate from reconstructed deterministic state.

## Stability guarantees

DigiEmu Core v1.x requires the following stability guarantees:

- same deterministic state -> same canonical JSON -> same SHA-256 hash
- different input key ordering -> same canonical JSON, if logical content is identical
- pretty-printed or compact JSON -> same canonical JSON, if logical content is identical
- replaying the same bundle input -> same reconstructed state hash
- tampered deterministic state -> different hash -> verification FAIL

Ordering constraints and determinism requirements are normative. See `docs/PUBLIC_API_README.md` and `docs/VERIFY_SPEC_v1.0.md` for related contracts.