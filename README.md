# DigiEmu Core

Deterministic Knowledge Infrastructure for AI Systems.

DigiEmu Core provides a deterministic snapshot and verification layer for knowledge states.  
It allows AI and software systems to store, replay, verify and sign knowledge states with full reproducibility.

The system ensures that knowledge snapshots can always be reconstructed and verified.

---

## Core Principles

- Deterministic knowledge snapshots
- Cryptographic integrity
- Replay verification
- Signed knowledge bundles
- Trusted identity layer
- Transportable bundles

---

## Specification

- [DigiEmu Core Specification v0.9](docs/DIGIEMU_CORE_SPEC_v0.9.md) — public review draft consolidating the current deterministic snapshot, hashing and replay verification contracts.
- [Spec Index v1.0](docs/SPEC_INDEX_v1.0.md) — index of normative DigiEmu Core documents.
- [Snapshot Hash v1.0](docs/SNAPSHOT_HASH_v1.0.md) — canonical JSON and SHA-256 snapshot identity.
- [Verify Spec v1.0](docs/VERIFY_SPEC_v1.0.md) — replay verification and PASS / FAIL result contract.

---

## Quick Start

### Build the CLI

```bash
go build -o digiemu ./cmd/digiemu
```

### Create a snapshot

```bash
digiemu snapshot file input.json
```

### Verify bundle

```bash
digiemu verify bundle snapshots/.../bundle.json
```

### Replay snapshot

```bash
digiemu replay bundle snapshots/.../bundle.json
```

### Verify replay determinism

```bash
digiemu verify replay snapshots/.../bundle.json
```

---

## Signature System

### Sign a bundle

```bash
digiemu sign bundle bundle.json
```

### Verify signature

```bash
digiemu verify signature bundle.json
```

---

## Identity System

### Show local identity

```bash
digiemu identity show
```

### Export identity

```bash
digiemu identity export <directory>
```

### Import trusted identity

```bash
digiemu identity import <directory>
```

### Show identity fingerprint

```bash
digiemu identity fingerprint
```

---

## Bundle Transport

### Export bundle

```bash
digiemu export bundle bundle.json <directory>
```

### Import bundle

```bash
digiemu import bundle <directory>
```

---

## End-to-End Verification Pipeline

DigiEmu Core guarantees that a knowledge snapshot can always be:

- recreated
- replayed
- verified
- cryptographically signed
- transported
- trusted via identity verification

---

## Architecture Overview

DigiEmu Core consists of:

- Snapshot engine
- Replay engine
- Verification layer
- Signature system
- Identity trust layer
- Bundle transport system

---

## Repository Structure

```text
cmd/digiemu
pkg/snapshot
pkg/replay
pkg/verify
pkg/claims
pkg/meaning
pkg/uncertainty
```

---

## Status

DigiEmu Core v1.0.0 is the first deterministic baseline release.

The system supports:

- deterministic snapshot creation
- replay verification
- cryptographic integrity validation
- identity-bound artifacts

Enterprise hardening and ecosystem integration are ongoing.

---

## Release

Current version: `v1.0.0`

See:

```text
https://github.com/DigiEmu/core/releases
```

---

## License

Business Source License 1.1.

See [LICENSE](LICENSE) for details.

---

## Security

- [Security Policy](SECURITY.md)
- [Determinism Exceptions](docs/security/DETERMINISM_EXCEPTIONS.md)