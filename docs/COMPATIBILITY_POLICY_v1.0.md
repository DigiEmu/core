# DigiEmu Core — Compatibility Policy v1.0

## Semantic Versioning

MAJOR.MINOR.PATCH

- MAJOR: breaking API or determinism contract change
- MINOR: backwards-compatible feature addition
- PATCH: bugfixes that do not affect canonical bytes or hashes

## Determinism is part of the contract

Any change affecting the following requires a MAJOR version bump:
- canonical encoding rules
- hash scope rules
- snapshot manifest/state formats
- replay semantics (state reconstruction)

## Golden vectors

Canonical encoding and hashing behavior is protected by golden test vectors.
If golden outputs change, the change must be treated as MAJOR.
