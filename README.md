# DigiEmu Core

Deterministic and verifiable knowledge infrastructure for AI systems.

DigiEmu Core provides a deterministic snapshot and verification layer for knowledge states.  
It enables AI and software systems to store, replay, verify, and sign knowledge states with full reproducibility.

The system ensures that knowledge snapshots can always be reconstructed and independently verified.

---------------------------------------------------------------------

## Core Principles

- Deterministic knowledge snapshots
- Cryptographic integrity
- Replay verification
- Signed knowledge bundles
- Trusted identity layer
- Transportable bundles

---------------------------------------------------------------------

## Quick Start

Build the CLI:

```bash
go build -o digiemu ./cmd/digiemu

Create a snapshot:

digiemu snapshot file input.json

Verify bundle:

digiemu verify bundle snapshots/.../bundle.json

Replay snapshot:

digiemu replay bundle snapshots/.../bundle.json

Verify replay determinism:

digiemu verify replay snapshots/.../bundle.json
Signature System

Sign a bundle:

digiemu sign bundle bundle.json

Verify signature:

digiemu verify signature bundle.json
Identity System

Show local identity:

digiemu identity show

Export identity:

digiemu identity export <directory>

Import trusted identity:

digiemu identity import <directory>

Identity fingerprint:

digiemu identity fingerprint
Bundle Transport

Export bundle:

digiemu export bundle bundle.json <directory>

Import bundle:

digiemu import bundle <directory>
End-to-End Verification Pipeline

DigiEmu Core guarantees that a knowledge snapshot can always be:

recreated
replayed
verified
cryptographically signed
transported
trusted via identity verification
Architecture Overview

DigiEmu Core consists of:

Snapshot engine
Replay engine
Verification layer
Signature system
Identity trust layer
Bundle transport system
Repository Structure
cmd/digiemu
pkg/snapshot
pkg/replay
pkg/verify
pkg/claims
pkg/meaning
pkg/uncertainty
Status

DigiEmu Core v1.0.0 is the first deterministic baseline release.

The system supports:

deterministic snapshot creation
replay verification
cryptographic integrity validation
identity-bound artifacts

Enterprise hardening and ecosystem integration are ongoing.

Release

Current version: v1.0.0

See:
https://github.com/DigiEmu/core/releases

License

Business Source License (BSL).
See LICENSE for details.

## 🔐Security

- [Security Policy](SECURITY.md)  
- [Determinism Exceptions](docs/security/DETERMINISM_EXCEPTIONS.md)