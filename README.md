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
- Artifact origin metadata
- Transportable bundles

---

## Boundary Note

DigiEmu Core verifies deterministic decision-state integrity through canonical JSON, SHA-256 snapshot hashing, and replay/verify checks. It does not certify agent identity, trustworthiness, authorization, action legitimacy, or legal compliance. These concerns belong to complementary external trust and attestation layers such as TBN.

## Specification

DigiEmu Core is moving toward a public standard structure for deterministic AI decision verification.

### Public review drafts

- [DigiEmu Core Specification v0.9](docs/DIGIEMU_CORE_SPEC_v0.9.md)
- [Test Vectors v0.9](docs/TEST_VECTORS_v0.9.md)
- [Negative Test Vectors v0.9](docs/NEGATIVE_TEST_VECTORS_v0.9.md)
- [Test Vector Manifest v0.9](docs/TEST_VECTOR_MANIFEST_v0.9.json)
- [Conformance v0.9](docs/CONFORMANCE_v0.9.md)
- [Conformance Declaration v0.9](docs/CONFORMANCE_DECLARATION_v0.9.md)
- [Conformance Declaration Schema v0.9](docs/CONFORMANCE_DECLARATION_SCHEMA_v0.9.json)
- [Verify Report Examples v0.9](docs/VERIFY_REPORT_EXAMPLES_v0.9.md)

### Specification index and implementation contracts

- [Spec Index v1.0](docs/SPEC_INDEX_v1.0.md)
- [CLI Contract v1.0](docs/CLI_CONTRACT_v1.0.md)
- [Snapshot Hash v1.0](docs/SNAPSHOT_HASH_v1.0.md)
- [Verify Spec v1.0](docs/VERIFY_SPEC_v1.0.md)
- [Verify Result Schema v1](docs/VERIFY_RESULT_SCHEMA_v1.json)
- [Verify Report Schema v0.9](docs/VERIFY_REPORT_SCHEMA_v0.9.json)
- [Snapshot Bundle v1.0](docs/SNAPSHOT_BUNDLE_v1.0.md)

---

### Public standard path

```text
Specification -> Test Vectors -> Negative Test Vectors -> Test Vector Manifest -> Conformance -> Conformance Declaration -> Conformance Declaration Schema -> Verify Report Examples -> Verify Report Schema
```

The specification explains the model.  
The test vectors make verification reproducible.  
The conformance document defines implementation requirements.  
The verify report examples define machine-readable verification outcomes.

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

## Artifact origin / bundle signer metadata

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
- Artifact origin / signer metadata utilities
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
- origin-bound or externally attested artifacts

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

## Core 2.0 Draft 1 — Partner-Testable Milestone

Status: draft / pre-release / not stable

Core 2.0 Draft 1 (tag example: `v2.0.0-draft.1`) is a partner-testable milestone.
It is intended for evaluation and feedback; it is NOT a stable release and does
not change any v1.0 behavior.

Short positioning
- DigiEmu Core provides deterministic knowledge infrastructure for AI systems.
- Core 2.0 Draft 1 is partner-testable and draft-only.
- v1.0 behavior remains unchanged and is the baseline for compatibility.

Core 2.0 capabilities (draft)
- Experimental conformance CLI
- Human-readable conformance output
- JSON conformance reports
- Conformance report schema
- OpenAPI contract draft
- Docker usage path (documentation)
- CI conformance checks

Quickstart commands (for evaluation)
```bash
go test ./...
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance

“DigiEmu Core by Bruno Baumgartner © 2026”
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance --json
```

Docker (optional evaluation)
```bash
docker build -t digiemu-core .
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance --json
```

Expected output
- `Conformance run summary: total=11 passed=11 failed=0`

Important notes
- Core 2.0 remains draft unless explicitly marked stable.
- Experimental commands remain under the `experimental` namespace.
- This is not a stable `v2.0.0` release.

Links
- [Conformance Quickstart](docs/CORE_2_CONFORMANCE_QUICKSTART.md)
- [Partner Integration Notes](docs/CORE_2_PARTNER_INTEGRATION_NOTES.md)
- [Docker Usage](docs/CORE_2_DOCKER_USAGE.md)
- [OpenAPI Draft](docs/CORE_2_OPENAPI_DRAFT.md)
- [Release Checklist](docs/CORE_2_RELEASE_CHECKLIST.md)
- [Tagging Plan](docs/CORE_2_TAGGING_PLAN.md)


