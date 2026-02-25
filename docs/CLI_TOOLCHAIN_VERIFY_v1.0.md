# CLI Toolchain — `digiemu verify` (Normative) v1.0

Status: NORMATIVE
Scope: Public toolchain contract (CLI surface), not internal implementation details.

## Purpose
`digiemu verify` is the canonical CLI entrypoint to verify a snapshot reference deterministically.
It MUST be stable across v1.x and MUST remain compatible with the `pkg/verify` and `pkg/snapshot` public contracts.

This document defines:
- Required CLI interface (flags/args)
- Exit codes
- Output format (JSON schema for machine parsing)
- Determinism expectations

## Command
### Synopsis (REQUIRED)
```
digiemu verify --ref <REF> [--data <DIR>] [--json]
```

### Arguments / Flags
- `--ref <REF>` (REQUIRED)
  - A snapshot reference string. The exact grammar MUST match the `pkg/snapshot` reference format.
- `--data <DIR>` (OPTIONAL; default `./data`)
  - Filesystem data directory used by the verifier implementation (adapters).
- `--fixture-root <DIR>` (OPTIONAL; default `data/test-fixtures`)
  - Optional fixture directory. The verifier will look for bundles under `<fixture-root>/snapshots/<ref>/snapshot.json` before falling back to `--data` unless `--prefer-data` is set.
- `--prefer-data` (OPTIONAL; default `false`)
  - If set, prefer `--data` paths over fixtures when both exist.
- `--json` (OPTIONAL; default `false`)
  - When set, output MUST be a single JSON object on stdout.

Implementations MAY add additional flags, but MUST NOT change the meaning of the required ones.

## Exit Codes (REQUIRED)
- `0` verification succeeded (`ok=true`)
- `1` verification failed (`ok=false`) with a valid result object on stdout
- `2` usage / input error (missing `--ref`, malformed flags); MAY print usage to stderr
- `>2` reserved for unexpected runtime errors (I/O, internal panic recovered, etc.)

## Output Format (REQUIRED)
When `--json` is enabled, stdout MUST be:
- exactly one JSON object (no extra lines), UTF-8

### JSON fields (REQUIRED)
- `ok` (boolean)
- `ref` (string) — the snapshot ref/hash used
- `expected` (string) — expected hash from the bundle
- `got` (string) — computed hash
- `hash_alg` (string) — e.g. `sha256(canonical_json_v1)`
- `canonical_scope` (string) — canonicalization scope description
- `errors` (array[string]) — attempted paths, IO errors, or validation messages
- `message` (string, optional; human-readable)

JSON MUST be a single object followed by a newline. Consumers should rely on the defined fields above.

### Example (SUCCESS)
```json
{"ok":true,"message":"verified","ref":{"hash":"...","algo":"sha256","format":"cbor","scope":"..." }}
```

### Example (FAILURE)
```json
{"ok":false,"message":"hash mismatch","ref":{"hash":"...","algo":"sha256","format":"cbor","scope":"..." }}
```

## Determinism Requirements
- Given the same `--ref` and the same referenced inputs, the result MUST be stable (same ok/message class).
- The verifier MUST NOT depend on wall-clock time for decision making.
- Any non-deterministic behavior MUST be surfaced as `ok=false` with a clear message.

## Compatibility Notes
- This CLI contract is part of the public toolchain surface.
- Internal implementations may change as long as they preserve:
  - flags/exit codes
  - output schema
  - determinism expectations
