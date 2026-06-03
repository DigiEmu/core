# Core 2.0 Conformance Report (JSON)

This document describes the machine-readable conformance report emitted by the
experimental conformance CLI when invoked with `--json`.

Report shape
------------

Top-level fields:

- `report_version` (string) — report format version, example: `core-2-conformance-report-v1`.
- `status` (string) — overall status: `PASS` when all cases structurally validated, `FAIL` otherwise.
- `total` (integer) — total number of cases discovered.
- `passed` (integer) — number of cases considered passed by the runner (structurally valid expectations).
- `failed` (integer) — number of cases with structural errors.
- `cases` (array) — per-case objects with fields:
  - `name` — case directory name
  - `case_passed` — boolean; true when expected_verify_result.json is structurally valid
  - `expected_result` — declared expected verify result (PASS/FAIL/ERROR)
  - `reason_code` — declared reason code
  - `error` — optional error message if structural validation failed

Schema
------

See `schemas/core_2_conformance_report.schema.json` for the JSON Schema of the report.

CLI usage
---------

Human-readable summary (default):

```
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance
Conformance run summary: total=11 passed=11 failed=0
```

Machine-readable JSON output:

```
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance --json
{ ... }
```

Notes
-----

- The `--json` flag is experimental and does not change the default human-readable output.
- The report schema is draft and intended for partner integration and ingestion.
