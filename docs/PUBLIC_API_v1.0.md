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

## Snapshot (pkg/snapshot)

**Purpose:** Stable public types to reference and describe a deterministic knowledge snapshot.

**Key types**
- `snapshot.Hash` — stable hex string (e.g. sha256)
- `snapshot.Ref` — `{ hash }` identifier (public reference)
- `snapshot.Manifest` — minimal public manifest shape:
	- `version` (string)
	- `hash` (Hash)
	- `inputs` ([]string)

**Rules**
- `Ref.Validate()` MUST fail if `hash` is empty.
- `Manifest` is the public stable shape; internal engines may carry richer state.

## Verify (pkg/verify)

**Purpose:** Stable interface for snapshot verification (internal implementations live under `internal/*`).

**Key types**
- `verify.Result` — `{ ok, message?, ref }`
- `verify.Verifier` — `Verify(ref snapshot.Ref) (Result, error)`

**Rules**
- Verifiers MUST return a `Result` with the same `ref` they verified.
- `Result.Message` is optional and should be used for human-readable failure context.

## Meaning (pkg/meaning)

**Purpose:** Stable public types for expressing meaning metadata and constraints on interpretation.

**Rules**
- Public types in `pkg/meaning` are stable and backwards-compatible within `v1.x`.
- Validation helpers (if present) must be deterministic and side-effect free.

## Claims (pkg/claims)

**Purpose:** Stable public types for representing verifiable claims and their identifiers.

**Rules**
- Claim identifiers must be stable strings.
- Validation helpers (if present) must be deterministic and side-effect free.

## Uncertainty (pkg/uncertainty)

**Purpose:** Stable public types for uncertainty as a first-class domain object.

**Rules**
- Uncertainty must be representable without requiring internal packages.
- Validation helpers (if present) must be deterministic and side-effect free.
