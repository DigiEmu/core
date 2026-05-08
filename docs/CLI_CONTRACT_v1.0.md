# DigiEmu Core — CLI Contract v1.0

Status: DRAFT  
Scope: `digiemu` CLI external contract

This document defines the normative, stable external contract of the `digiemu` CLI.

It is intentionally minimal: it specifies what is guaranteed, what is not guaranteed, and what must remain stable across v1.x releases.

---

## 1. Scope

This contract applies to the `digiemu` CLI shipped from:

```text
./cmd/digiemu
```

The CLI is considered a tool interface and therefore MUST:

- be scriptable
- be deterministic for deterministic commands
- provide stable JSON output where defined
- provide governance-stable exit codes where defined

This contract does not define internal Go package APIs.

---

## 2. Stability Levels

The following MUST remain stable in v1.x:

- command names and subcommand structure listed in this document
- flag names listed in this document, with behavior as specified
- stable JSON output fields for commands that define JSON output
- exit code semantics for commands that define exit codes
- determinism rules for deterministic commands

The following MAY change in v1.x:

- human-readable non-JSON output formatting
- help text wording
- internal trace details beyond the normative trace rules
- performance characteristics

The following are NOT guaranteed:

- undocumented flags
- undocumented command aliases
- undocumented output fields
- internal package behavior not covered by this contract

---

## 3. Supported Commands

The following commands are part of the documented CLI surface:

```text
digiemu verify
digiemu replay
digiemu unit create
digiemu version create
digiemu audit verify
digiemu audit tail
digiemu export unit
digiemu meaning set
digiemu meaning show
digiemu claim set
digiemu claim show
digiemu uncertainty set
digiemu uncertainty show
```

The following command is available but considered experimental unless otherwise specified:

```text
digiemu serve
```

Notes:

- This contract focuses on `verify` and `replay` as the deterministic v1.0 tool core.
- Other commands remain supported, but their detailed schemas may be governed by their own documents or future contracts.

---

## 4. `verify` Contract

### 4.1 Command

```text
digiemu verify
```

### 4.2 Stable Flags

```text
--ref <ref>
--bundle <path>
--data <dir>
--fixture-root <dir>
--prefer-data
--json
--write-expected
--strict
```

Flag semantics:

- `--ref <ref>` identifies the snapshot reference. It is required unless `--bundle` is set.
- `--bundle <path>` points directly to a bundle root. When set, it is the single source of truth.
- `--data <dir>` defines the runtime data root. Default: `./data`.
- `--fixture-root <dir>` defines the fixture root. Default: `data/test-fixtures`.
- `--prefer-data` changes resolver order when `--bundle` is not set.
- `--json` emits stable JSON output.
- `--write-expected` enables governed expected-hash write behavior. It MUST NOT overwrite a non-placeholder expected hash.
- `--strict` is legacy-compatible and MUST NOT override documented exit code semantics.

### 4.3 Bundle vs Ref Resolution

If `--bundle` is set:

- the provided bundle root MUST be used as the source of truth
- resolver flags such as `--data`, `--fixture-root`, and `--prefer-data` MUST NOT affect which bundle is used
- if `--ref` is missing, the reference SHOULD be derived from the bundle directory name when possible

If `--bundle` is not set:

- resolver behavior MUST follow `docs/SNAPSHOT_BUNDLE_v1.0.md`

### 4.4 JSON Output

When `--json` is set, `verify` MUST output JSON conforming to:

```text
docs/VERIFY_RESULT_SCHEMA_v1.json
```

Required high-level fields include:

```text
ok
ref
expected
got
hash_alg
canonical_scope
trace
message
```

If write policy is in play, the following fields MUST be present:

```text
wrote_expected
write_blocked
write_reason
```

### 4.5 Determinism

For a fixed bundle, `verify --json` MUST be deterministic in:

- computed hash values
- exit code
- JSON field values
- trace ordering

Exception:

- absolute path prefixes may differ across machines

Path determinism rule:

- trace entries MUST preserve stable relative ordering
- trace entries SHOULD preserve stable suffixes
- absolute workspace prefixes MAY differ across machines

### 4.6 Exit Codes

Exit codes are governed by:

```text
docs/CLI_VERIFY_v1.0.md
```

This contract references that document as normative for `verify` exit code semantics.

---

## 5. `replay` Contract

### 5.1 Command

```text
digiemu replay
```

### 5.2 Stable Flags

```text
--bundle <path>
--json
```

Flag semantics:

- `--bundle <path>` points to the bundle root and is required.
- `--json` emits stable JSON output.

### 5.3 JSON Output

When `--json` is set, `replay` MUST emit JSON containing:

```text
snapshot
claims
trace
```

Minimum expectations:

- `snapshot` MUST be an object containing at least `ref` and `state` when available.
- `claims` MUST be an array and MAY be empty.
- `trace` MUST be an array of strings.

### 5.4 Determinism

For a fixed bundle, `replay --json` MUST be byte-stable across repeated runs on the same machine.

Across machines:

- absolute path prefixes may differ
- trace ordering MUST remain deterministic
- the `used:` marker rule MUST hold

---

## 6. Trace Rules

For deterministic commands that emit trace, including `verify` and `replay`:

- trace MUST include the files read, at minimum `snapshot.json`
- claim files MUST be included in trace if present and read
- trace MUST end with exactly one final `used:<bundle-root>` entry
- the `used:` marker MUST be the last trace entry
- trace ordering MUST be deterministic

---

## 7. Compatibility Rules

Minor releases in v1.x MUST NOT break:

- documented command names
- documented flag names
- JSON schema meaning
- exit code semantics
- deterministic behavior for deterministic commands

Minor releases in v1.x MAY add:

- new optional JSON fields
- new commands
- new flags
- additional trace entries, provided normative trace ordering remains stable

Existing documented fields MUST retain their meaning.

---

## 8. Non-Goals

This contract does not guarantee:

- stable formatting of non-JSON output
- stable ordering of map keys in non-canonical contexts
- stable performance characteristics
- stable memory use
- stable internal package names
- stable internal Go types
- production stability of experimental commands such as `digiemu serve`

---

## 9. References

Normative:

```text
docs/CLI_VERIFY_v1.0.md
docs/SNAPSHOT_BUNDLE_v1.0.md
docs/VERIFY_RESULT_SCHEMA_v1.json
docs/VERIFY_SPEC_v1.0.md
docs/SNAPSHOT_HASH_v1.0.md
```