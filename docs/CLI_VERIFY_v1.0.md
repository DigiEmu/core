# digiemu verify — CLI (v1.0)

Usage:

```
digiemu verify --ref <hash> [--data <data-dir>] [--format text|json] [--strict]
```

Behavior:

- Reads snapshot bundle from `<data>/snapshots/<ref>/snapshot.json`.
- Computes `hash_v1 = SHA-256(canonical_json(state))` using the internal canonicalizer.
- Compares `hash_v1` to `expected_hash_v1` in the bundle.

Exit codes:

- `0` = verification succeeded (hash match)
- `2` = verification failed (mismatch or bundle error)
- `1` = unexpected runtime error

Output formats:

- `--format=json` — prints a single JSON object (the `pkg/verify.Result` shape).
- `--format=text` — prints a short text line: `OK\t<ref>` or `FAIL\t<ref>\t<msg>`.

Link: see `docs/SNAPSHOT_BUNDLE_v1.0.md` for bundle layout and `docs/SNAPSHOT_HASH_v1.0.md` for the canonical hashing spec.
CLI — verify (v1.0 stub)
This command is a deterministic placeholder that validates the snapshot ref and returns OK for non-empty hashes.

Usage
  digiemu verify --ref <snapshot-hash> [--format text|json] [--strict]

Examples
  digiemu verify --ref 0123abcd --format text
  digiemu verify --ref 0123abcd --format json
  digiemu verify --ref "" --strict

Exit behavior
- Non-strict (default):
  - returns 0 if the command ran and produced output (even if verification is FAIL due to invalid ref handling detail)
  - returns 1 only for command/runtime errors (bad flags etc.)
- Strict:
  - returns 1 if verification is not OK, or if Verify returns an error

Notes
- Phase 6 will replace the stub with a real engine-backed verifier.
- The public API types live in pkg/snapshot and pkg/verify.
