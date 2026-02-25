# digiemu verify — CLI (v1.0)

Usage:

```
digiemu verify --ref <hash> [--data <data-dir>] [--fixture-root <dir>] [--prefer-data] [--bundle <path>] [--json] [--write-expected] [--strict]
```

Behavior:

- Reads snapshot bundle from `<fixture-root>/snapshots/<ref>/` or `<data>/snapshots/<ref>/`.
- Alternatively pass `--bundle <path>` to load a bundle root directly (bypass ref lookup).
- Computes `hash_v1 = SHA-256(canonical_json(state))` using the internal canonicalizer.
- Compares `hash_v1` to `expected_hash_v1` in the bundle.

## Bundle vs ref (normative)

- `--bundle` makes `--ref` optional. If `--bundle` is set, it is the single source of truth for which files are read.
- When `--bundle` is provided, resolution flags (`--data`, `--fixture-root`, `--prefer-data`) MUST NOT influence file selection.
- When `--bundle` is not provided, `--ref` is required and resolver order follows `--prefer-data` as defined in `docs/SNAPSHOT_BUNDLE_v1.0.md`.
- Trace MUST end with `"used:<root>"`.

Minimal examples (expected exit codes shown):

```bash
# verify ok => 0
digiemu verify --bundle ./data/test-fixtures/snapshots/demo --json

# mismatch => 2 (after setting expected_hash_v1 to an incorrect value)
digiemu verify --bundle ./path/to/temp/snapshots/demo --json

# write blocked existing expected (with --write-expected) => 3
digiemu verify --bundle ./data/test-fixtures/snapshots/demo --write-expected --json
```

## Exit codes (normative)

The command `digiemu verify` SHALL return the following exit codes.

Notes:

- `--strict` does not override these exit codes (exit codes are deterministic regardless of `--strict`).

| Code | Meaning | When it occurs | Notes |
| ---: | --- | --- | --- |
| 0 | OK | Verification succeeds (`ok=true`). | Also ok when `--write-expected` is not set. |
| 2 | Hash mismatch | `expected != got`. | Verification failure. |
| 3 | Write blocked by governance | Only when `--write-expected` is set AND `expected_hash_v1` is already present (non-placeholder), i.e. would require overwrite. | This code may occur even if verification itself would otherwise be OK. |
| 4 | Invalid input | Missing snapshot, invalid JSON, invalid hash, invalid bundle path, IO error, or a write-blocked reason other than `existing_expected_present`. | Includes write-blocked reasons such as `snapshot_not_found`, `snapshot_invalid_json`, `invalid_hash`, `io_error`. |
| 5 | Internal / unexpected error | Unexpected runtime error (not classified as an input/write-policy error). | Implementation may exit before emitting JSON output. |

### Minimal examples (PowerShell)

Code 0 (OK):

```powershell
./digiemu.exe verify --bundle .\data\test-fixtures\snapshots\demo --json
echo $LASTEXITCODE
```

Code 2 (hash mismatch):

```powershell
$tmp = Join-Path $env:TEMP ("digiemu-mismatch-" + [guid]::NewGuid())
Copy-Item .\data\test-fixtures\snapshots\demo (Join-Path $tmp "snapshots\demo") -Recurse
$p = Join-Path $tmp "snapshots\demo\snapshot.json"
$o = Get-Content $p -Raw | ConvertFrom-Json
$o.expected_hash_v1 = ("0" * 64)
$o | ConvertTo-Json -Depth 100 | Out-File $p -Encoding utf8
./digiemu.exe verify --bundle (Join-Path $tmp "snapshots\demo") --json
echo $LASTEXITCODE
```

Code 3 (write blocked by governance):

```powershell
./digiemu.exe verify --bundle .\data\test-fixtures\snapshots\demo --write-expected --json
echo $LASTEXITCODE
```

Code 4 (invalid input; snapshot not found):

```powershell
./digiemu.exe verify --bundle .\data\test-fixtures\snapshots\does-not-exist --json
echo $LASTEXITCODE
```

Code 5 (internal/unexpected error example; force write failure):

```powershell
$tmp = Join-Path $env:TEMP ("digiemu-internal-" + [guid]::NewGuid())
Copy-Item .\data\test-fixtures\snapshots\demo (Join-Path $tmp "snapshots\demo") -Recurse
$p = Join-Path $tmp "snapshots\demo\snapshot.json"
$o = Get-Content $p -Raw | ConvertFrom-Json
$o.expected_hash_v1 = "REPLACE_AFTER_COMPUTE"
$o | ConvertTo-Json -Depth 100 | Out-File $p -Encoding utf8
attrib +R $p
./digiemu.exe verify --bundle (Join-Path $tmp "snapshots\demo") --write-expected --json
echo $LASTEXITCODE
```

Exit codes:

- `0` = verification succeeded (hash match)
- `2` = hash mismatch (`expected != got`)
- `3` = write blocked by governance (existing expected present; only with `--write-expected`)
- `4` = invalid input/bundle/snapshot/json/hash
- `5` = internal/unexpected error

Output formats:

- `--format=json` or `--json` — prints a single JSON object (the `pkg/verify.Result` shape).
- `--format=text` — prints a short text line: `OK\t<ref>` or `FAIL\t<ref>\t<msg>`.

Link: see `docs/SNAPSHOT_BUNDLE_v1.0.md` for bundle layout and `docs/SNAPSHOT_HASH_v1.0.md` for the canonical hashing spec.
CLI — verify (v1.0 stub)
This command is a deterministic placeholder that validates the snapshot ref and returns OK for non-empty hashes.

Usage
  digiemu verify --ref <snapshot-hash> [--format text|json] [--strict]

Examples
  digiemu verify --ref 0123abcd --format text
  digiemu verify --ref 0123abcd --format json
  digiemu verify --ref "" --strict

Exit behavior
- Non-strict (default):
  - returns 0 if the command ran and produced output (even if verification is FAIL due to invalid ref handling detail)
  - returns 1 only for command/runtime errors (bad flags etc.)
- Strict:
  - returns 1 if verification is not OK, or if Verify returns an error

Notes
- Phase 6 will replace the stub with a real engine-backed verifier.
- The public API types live in pkg/snapshot and pkg/verify.

## Windows / PowerShell note: fc

In Windows PowerShell, `fc` is commonly an alias for `Format-Custom`, not a file comparison tool.

For file compare, use:

- `fc.exe file1 file2`

Example: deterministic replay compare (PowerShell):

```powershell
./digiemu.exe replay --bundle .\data\test-fixtures\snapshots\demo --json | Out-File .\tmp1.json -Encoding utf8
./digiemu.exe replay --bundle .\data\test-fixtures\snapshots\demo --json | Out-File .\tmp2.json -Encoding utf8
fc.exe .\tmp1.json .\tmp2.json
```
