# Snapshot Bundle (MVP) — v1.0

This is a minimal on-disk snapshot bundle format used to make `digiemu verify` real.
It is not the final long-term bundle schema.

## Layout

`<data>/snapshots/<ref>/snapshot.json`

Bundle layout (v1):

- Root directory: either `<fixture-root>/snapshots/<ref>/` or `<data>/snapshots/<ref>/` (fixture preferred unless `--prefer-data`)
- Required file: `snapshot.json` (must contain `expected_hash_v1` and may include a `state` object for compatibility)
- Optional directories: `units/`, `versions/`, `claims/`, `meaning/`, `uncertainty/`
	- If present, the verifier will read all `*.json` files from these directories, sorting filenames lexicographically before processing.
	- Each file is included verbatim (BOM-stripped) as a JSON element in the assembled state arrays.

Verification state assembly:

- The assembled state used for hashing is a JSON object:

```json
{
	"snapshot": <raw snapshot.json content>,
	"units": [ <raw file1>, <raw file2>, ... ],
	"versions": [ ... ],
	"claims": [ ... ],
	"meaning": [ ... ],
	"uncertainty": [ ... ]
}
```

- Hash computation: `hash_v1 = SHA-256(canonical_json_v1(assembled_state))`
- `canonical_scope` for verifier output: `canonical_json_v1`
- JSON files are decoded in a BOM-safe manner (leading UTF-8 BOM is stripped before parsing)
 - When assembling the state for hashing, files from optional directories are included in lexicographic filename order.
 - Filenames are considered opaque; sorting ensures determinism even if files were written in different orders.

## snapshot.json

Required fields:

- `ref` (string) — optional (CLI ref is accepted if omitted)
- `expected_hash_v1` (string) — **required**
- `state` (object) — **required**

The verifier computes:

- `hash_v1 = SHA-256(canonical_json(state))`

and compares it to `expected_hash_v1`.

## Determinism

Canonical JSON rules:
- map keys are sorted
- arrays keep order
- numbers/strings are encoded consistently

For the exact canonical encoding rules see `internal/canonicaljson` and
`docs/SNAPSHOT_HASH_v1.0.md`.

## Bundle root resolution (normative)

### BundleV1 root definition

A BundleV1 root SHALL be a directory of the form:

- `snapshots/<ref>/`
	- `snapshot.json` (required)
	- `claims/*.json` (optional)
	- other optional directories (e.g. `units/`, `versions/`, `meaning/`, `uncertainty/`) (optional)

The verifier/replay SHALL treat `snapshots/<ref>/snapshot.json` as the required entrypoint.

### Root resolution order

1) If `--bundle` is provided:
	 - It MUST be treated as the bundle root.
	 - The verifier/replay MUST NOT perform any ref-based lookup.
	 - The CLI MAY derive `--ref` from the directory name, but this MUST NOT change which files are read.
	 - The trace MUST end with `"used:<bundle_root>"`.

2) Else (no `--bundle`):
	 - `--ref` is required.
	 - The tool MUST attempt to find a bundle root in:
		 a) `--data/snapshots/<ref>`
		 b) `--fixture-root/snapshots/<ref>`
		 but the order MUST follow `--prefer-data`:
		 - if `--prefer-data` is true: try data first, then fixture-root
		 - else: try fixture-root first, then data
	 - If neither exists: the tool MUST fail deterministically (snapshot/bundle not found).

### Trace norms

- Trace MUST include every file path read (at minimum: `snapshot.json` and each loaded claim file).
- Trace MUST include exactly one final `"used:<root>"` marker indicating the selected root.
- Trace order SHOULD be stable and deterministic.

### Minimal examples

```bash
digiemu verify --bundle ./data/test-fixtures/snapshots/demo --json
digiemu verify --ref demo --fixture-root ./data/test-fixtures --data ./data --json
digiemu replay --bundle ./data/test-fixtures/snapshots/demo --json
```
