# Snapshot Hash v1.0 — Canonical JSON

Status: NORMATIVE

This document defines the canonical snapshot hashing algorithm used by DigiEmu Core v1.x.

Scope: canonical_utf8_without_sha256_comment_line, SHA-256, uppercase hex.

Rules (summary):
- Build an explicit replay state object that contains only the stable kernel state (units, versions, meaning, claims, uncertainty).
- Serialize using `internal/canonicaljson.Marshal` which sorts map keys, preserves array order, and emits compact JSON without extra whitespace.
- Compute SHA-256 over the UTF-8 bytes of the canonical JSON.
- Encode the digest as uppercase hexadecimal and record it where specified (e.g. `docs/GENESIS_ANCHOR_v1.0.yaml`), excluding the sha256 comment line from canonical bytes.

Ordering constraints and determinism requirements are normative — see `docs/PUBLIC_API_README.md` for links.
