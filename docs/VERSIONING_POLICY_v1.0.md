# DigiEmu Core — VERSIONING_POLICY_v1.0
**Status:** DRAFT (Phase 5)  
**Scope:** CLI contract + file formats + verify/replay determinism surface  
**Applies to:** `digiemu-core` repository and the `digiemu` CLI

## 1. Purpose (normative)
This document defines the versioning and compatibility policy for DigiEmu Core and its CLI.
The goal is to provide predictable upgrade paths and stable external contracts for tooling,
integrations, and auditors.

## 2. Definitions (normative)

### 2.1. Semantic versioning model
DigiEmu Core uses **MAJOR.MINOR.PATCH** versioning for **public releases**.

- **MAJOR**: breaking changes to any normative contract defined below.
- **MINOR**: backward-compatible feature additions or expansions.
- **PATCH**: backward-compatible bug fixes and internal refactors.

### 2.2. “Public contract” (normative)
A public contract includes, but is not limited to:

1) **CLI JSON result schemas**
   - `verify --json` output fields and their semantics
   - `replay --json` output fields and their semantics

2) **CLI exit code governance**
   - normative exit codes and their meanings

3) **Bundle / file format rules**
   - snapshot bundle layout (`snapshots/<ref>/snapshot.json`, etc.)
   - canonical hashing scope rules (e.g. `canonical_json_v1_excluding_expected_hash_v1`)

4) **Determinism rules**
   - “same inputs => byte-identical stable JSON output” requirements where specified

If a change modifies any of the above in an incompatible way, it is a breaking change.

## 3. Breaking change policy (normative)

### 3.1. What counts as breaking (MAJOR bump required)
Any of the following MUST trigger a MAJOR bump:

- Removing or renaming JSON fields in CLI outputs (verify/replay), or changing their type
- Changing exit code meanings or mappings
- Changing canonical hashing scope or algorithm identifiers for existing versions
- Changing bundle resolution rules in a way that changes which bundle is selected for the same flags
- Changing determinism guarantees (e.g., making JSON output ordering unstable)

### 3.2. What is NOT breaking (allowed in MINOR/PATCH)
The following are NOT breaking:

- Adding new JSON fields (MUST be additive; existing fields unchanged)
- Adding new subcommands or flags (without changing existing behavior)
- Internal refactors that preserve outputs and invariants
- Performance improvements
- Additional validation that rejects inputs previously undefined/invalid by spec (tightening within spec)

## 4. Deprecation policy (normative)

### 4.1. Deprecation signals
When something is deprecated:
- It MUST be documented in release notes.
- CLI behavior SHOULD emit a stable warning string on stderr (if applicable),
  but MUST NOT change exit codes unless explicitly specified.

### 4.2. Deprecation window
- Deprecations MUST remain supported for at least **one MINOR release**.
- Removal of deprecated behavior MUST be done only in a **MAJOR** release.

## 5. Compatibility promise (normative)

### 5.1. v1.x promise
During v1.x:
- No breaking changes to public contracts (Section 2.2).
- Additive improvements only (new optional fields, new commands).
- Exit codes remain stable.

### 5.2. Spec versions vs implementation versions
If a document is labeled `*_v1.0` (spec), it defines a contract.
Implementation versions MUST NOT violate finalized spec documents without a MAJOR bump.

## 6. Release discipline (normative)

- Every release MUST include a clear changelog entry.
- If a change touches a public contract, it MUST include:
  - updated docs
  - tests that lock the behavior (integration tests when appropriate)

## 7. References (informative)
- `docs/CLI_CONTRACT_v1.0.md`
- `docs/CLI_VERIFY_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
