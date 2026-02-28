# BUNDLE_LAYOUT_v1

Status: **Normative** (v1)

This document defines the required and optional on-disk layout for a **Snapshot Bundle** ("bundle") consumed by `digiemu verify --bundle <path>`.

## 1. Scope

This contract defines the filesystem layout and portability rules for bundles checked into repositories or transported between machines.

A bundle is a directory tree that MUST contain one or more snapshot roots under `snapshots/<ref>/`.

## 2. Directory layout (normative)

A bundle root directory MUST have the following structure:

```
<bundle-root>/
  snapshots/
    <ref>/
      snapshot.json
      hashes/                     (optional)
        *.sha256                  (optional)
      schemas/                    (optional)
        <schema-name>_v<ver>.json  (optional)
```

### 2.1 Required

- A bundle MUST contain a `snapshots/` directory.
- Each snapshot root MUST be located at `snapshots/<ref>/`.
- Each `snapshots/<ref>/` directory MUST contain `snapshot.json`.

### 2.2 Optional

- `snapshot.json` MAY include an `expected_hash_v1` field.
  - If present, it MUST be a UTF-8 text value.
- `hashes/` MAY be present at `snapshots/<ref>/hashes/`.
  - When present, it MAY contain additional `.sha256` files.
- `schemas/` MAY be present at `snapshots/<ref>/schemas/`.
  - When present, it MUST contain **locked schema copies** used to interpret the bundle.

## 3. Path portability rules (normative)

- Bundle contents MUST be usable via **relative paths only**.
  - Tools MUST NOT require absolute paths that embed a machine-specific prefix.
- Bundles MUST be safe to copy as plain files.
  - A bundle MUST NOT require symbolic links to be correct.

## 4. Filename rules (normative)

- Filenames inside a bundle SHOULD be deterministic and stable across platforms.
- Filenames inside a bundle SHOULD use **lowercase letters**, **digits**, and **underscores** (`[a-z0-9_]`).
- Tools consuming bundles MUST NOT depend on case-sensitive filesystem behavior.

## 5. Schema locking (normative)

To prevent silent/accidental schema drift:

- If a bundle includes `schemas/`, the files in `schemas/` MUST be treated as immutable, versioned artifacts.
- Locked schema filenames MUST be versioned (example: `verify_result_schema_v1.json`).
- A repository SHOULD additionally pin schema content via a committed hash (for example a `.sha256` lock file) and verify it in CI.

## 6. Determinism expectations (normative)

- Bundles MUST NOT require timestamps, local environment variables, or network access to verify.
- `snapshot.json` MUST be valid JSON.
- Verification results MUST be deterministic for a given bundle.

## 7. Examples

### 7.1 Minimal bundle

```
examples/bundles/demo_ok_bundle_v1/
  snapshots/
    snapshot_demo_v1/
      snapshot.json
```

### 7.2 Bundle with schema lock

```
examples/bundles/demo_ok_bundle_v1/
  snapshots/
    snapshot_demo_v1/
      snapshot.json
      expected_hash_v1
      schemas/
        verify_result_schema_v1.json
```
