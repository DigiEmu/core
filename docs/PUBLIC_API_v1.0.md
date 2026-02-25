# DigiEmu Core — Public API v1.0

Status: Draft (API Freeze Candidate)

## Scope

The stable public API surface for DigiEmu Core v1.0 is limited to:

- pkg/meaning
- pkg/claims
- pkg/uncertainty
- pkg/snapshot
- pkg/verify
- pkg/core (optional facade)

All other packages (including internal/*) are implementation details and may change without notice.

## Stability guarantee

DigiEmu Core follows Semantic Versioning (SemVer):

- v1.x.x: no breaking changes to the public API
- v2.0.0: breaking changes allowed

A breaking change includes:
- removal or signature change of exported identifiers
- changes to canonicalization behavior
- changes to hash scope behavior
- changes to snapshot artifact structure

## Deterministic contract

Public APIs producing canonical or hashed output must:
- be platform-independent
- not rely on ambient time
- not rely on OS-specific paths
- produce byte-identical output for identical inputs
