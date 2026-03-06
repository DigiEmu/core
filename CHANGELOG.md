# Changelog

## v0.1.0-core-mvp

DigiEmu Core MVP completed.

### Core capabilities

- deterministic snapshot engine
- bundle manifest generation
- bundle verification
- replay bundle command
- replay determinism verification
- snapshot hash integrity verification
- replay hash persistence
- bundle export/import foundation
- bundle signature layer
- trusted identity layer
- identity export/import
- identity fingerprint

### Result

DigiEmu Core now provides deterministic, portable, replayable, verifiable and signable knowledge artifacts.

A bundle can now be:

- created
- verified
- replayed
- replay-verified
- exported
- imported
- signed
- signature-verified
- bound to a trusted identity
- referenced by identity fingerprint

---

## Historical development notes

### v0.5.0 — Uncertainty Minimum v0

- Added optional version-scoped uncertainty metadata
- Added uncertainty hash persistence and verification
- Added CLI and HTTP endpoints for uncertainty
- Added strict-hash verification support for uncertainty sidecars

### v0.3.0 — Meaning Layer v1

- Added optional version-scoped meaning metadata
- Added meaning hash persistence and verification
- Added CLI and HTTP endpoints for meaning
- Added strict-hash verification support for meaning sidecars

### v0.2.9 — Audit & Snapshot Stabilization

- Deterministic snapshot hashing
- Hardened audit and snapshot pipeline
- Added export CLI command and audit verification tooling

### v0.1.0

- Kernel: Units & Versions with validation
- Ports & DTO contracts
- In-memory + FS JSON repository
- HTTP API with consistent error format
- CLI and HTTP API entrypoints
- CI workflow with formatting and tests