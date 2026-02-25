# Snapshot Bundle (MVP) — v1.0

This is a minimal on-disk snapshot bundle format used to make `digiemu verify` real.
It is not the final long-term bundle schema.

## Layout

`<data>/snapshots/<ref>/snapshot.json`

## snapshot.json

Required fields:

- `ref` (string) — optional (CLI ref is accepted if omitted)
- `expected_hash_v1` (string) — **required**
- `state` (object) — **required**

The verifier computes:

- `hash_v1 = SHA-256(canonical_json(state))`

and compares it to `expected_hash_v1`.

## Determinism

Canonical JSON rules:
- map keys are sorted
- arrays keep order
- numbers/strings are encoded consistently

For the exact canonical encoding rules see `internal/canonicaljson` and
`docs/SNAPSHOT_HASH_v1.0.md`.
