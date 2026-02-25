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
