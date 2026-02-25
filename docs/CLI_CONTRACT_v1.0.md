# DigiEmu Core — CLI_CONTRACT_v1.0
Status: DRAFT (Phase 5)

This document defines the **normative, stable external contract** of the `digiemu` CLI.
It is intentionally minimal: it specifies what is guaranteed, what is not, and what must remain stable.

---
## 1. Scope (normative)

This contract applies to the `digiemu` CLI shipped from `./cmd/digiemu`.

The CLI is considered a **tool interface** and therefore MUST:
- be scriptable
- be deterministic (for deterministic commands)
- provide stable JSON output (where defined)
- provide governance-stable exit codes (where defined)

This contract does **not** define internal package APIs.

---
## 2. Stability levels (normative)

**MUST remain stable in v1.x:**
- command names and subcommand structure listed in this document
- flag names listed in this document (behavior as specified)
- stable JSON output fields for commands that define JSON output
- exit code semantics for commands that define exit codes
- determinism rules for deterministic commands

**MAY change in v1.x:**
- human-readable (non-JSON) output formatting
- help text wording
- internal trace details beyond the normative rules
- performance characteristics

**NOT guaranteed:**
- any undocumented flag, command alias, or output field
- internal package behavior not covered by this contract

---
## 3. Supported commands (normative)

The following commands are part of the stable CLI surface:

- `digiemu verify`
- `digiemu replay`
- `digiemu unit create`
- `digiemu version create`
- `digiemu audit verify`
- `digiemu audit tail`
- `digiemu export unit`
- `digiemu serve`
- `digiemu meaning set`, `digiemu meaning show`
- `digiemu claim set`, `digiemu claim show`
- `digiemu uncertainty set`, `digiemu uncertainty show`

Notes:
- This contract focuses on `verify` and `replay` as the Phase 4/5 deterministic tool core.
- Other commands remain supported, but their detailed schemas are governed by their own docs/specs.

---
## 4. `verify` contract (normative)

### 4.1 Command
`digiemu verify`

### 4.2 Flags (stable)
- `--ref <ref>` (required unless `--bundle` is set)
- `--bundle <path>` (optional; when set, it is the single source of truth)
- `--data <dir>` (default `./data`)
- `--fixture-root <dir>` (default `data/test-fixtures`)
- `--prefer-data` (resolver order when `--bundle` is not set)
- `--json` (emit stable JSON)
- `--write-expected` (governed write policy; never overwrite)
- `--strict` (legacy; does not override exit codes)

### 4.3 Bundle vs ref (normative)
- If `--bundle` is set, the bundle root MUST be used as the source of truth.
- If `--bundle` is set and `--ref` is missing, the ref MUST be derived from the bundle directory name (`snapshots/<ref>`).
- When `--bundle` is set, resolver flags (`--data`, `--fixture-root`, `--prefer-data`) MUST NOT affect which bundle is used.

When `--bundle` is not set:
- Resolver behavior MUST follow `docs/SNAPSHOT_BUNDLE_v1.0.md` (Bundle root resolution section).

### 4.4 JSON output (normative)
When `--json` is set, `verify` MUST output JSON that conforms to:
- `docs/VERIFY_RESULT_SCHEMA_v1.json`

Required high-level properties (minimum):
- `ok` (boolean)
- `ref` (string)
- `expected` (string or empty)
- `got` (string or empty)
- `hash_alg` (string)
- `canonical_scope` (string)
- `trace` (array of strings)
- `message` (string)

If write policy is in play, the following MUST be present:
- `wrote_expected` (boolean)
- `write_blocked` (boolean)
- `write_reason` (string enum)

### 4.5 Determinism (normative)
For a fixed bundle, `verify --json` MUST be deterministic in:
- computed hash values
- exit code
- JSON field values (including `trace` ordering) EXCEPT for absolute path prefixes on different machines

Path determinism rule:
- Trace entries MUST be stable relative ordering and stable suffixes.
- Absolute prefix differences (workspace paths) are allowed across machines.

### 4.6 Exit codes (normative)
Exit codes are governed by `docs/CLI_VERIFY_v1.0.md` (“Exit codes (normative)” section).
This contract references that document as normative for exit code semantics.

---
## 5. `replay` contract (normative)

### 5.1 Command
`digiemu replay`

### 5.2 Flags (stable)
- `--bundle <path>` (required)
- `--json` (emit stable JSON)

### 5.3 JSON output (normative)
When `--json` is set, `replay` MUST emit JSON with:
- `snapshot` object (containing at minimum `ref` and `state`)
- `claims` array (may be empty)
- `trace` array (strings)

### 5.4 Determinism (normative)
For a fixed bundle, `replay --json` MUST be byte-identical across repeated runs **on the same machine**.

Across machines:
- Absolute path prefixes may differ, but the trace ordering and the `used:` marker rule MUST hold.

---
## 6. Trace rules (normative)

For deterministic commands that emit trace (`verify`, `replay`):
- Trace MUST include the files read (at least snapshot.json, and claims files if present).
- Trace MUST end with **exactly one** final entry: `used:<bundle-root>`
- The `used:` marker MUST be the last trace entry.

---
## 7. Compatibility rules (normative)

- Minor releases (v1.x) MUST NOT break JSON schema or exit code semantics.
- New optional JSON fields MAY be added, but existing fields MUST retain meaning.
- New commands MAY be added, but existing commands MUST remain available.

---
## 8. Non-goals (explicit)

This contract does not guarantee:
- stable formatting of non-JSON output
- stable ordering of map keys in *non-canonical* contexts
- stable performance or memory use
- stable internal package names/types

---
## 9. References

Normative:
- `docs/CLI_VERIFY_v1.0.md`
- `docs/SNAPSHOT_BUNDLE_v1.0.md`
- `docs/VERIFY_RESULT_SCHEMA_v1.json`
- `docs/VERIFY_SPEC_v1.0.md`
